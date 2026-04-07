package pipeline

import (
	"context"

	"sync"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/dispatcher"
	"github.com/Anshuman-02905/chronostream/internal/dlq"
	"github.com/Anshuman-02905/chronostream/internal/engine"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
	"github.com/Anshuman-02905/chronostream/internal/scheduler"
	"github.com/Anshuman-02905/chronostream/internal/sequence"
	"github.com/Anshuman-02905/chronostream/internal/transport"
	"github.com/Anshuman-02905/chronostream/internal/user"
)

type FrequencyPipeline struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc

	Freq              event.Frequency
	Scheduler         scheduler.Scheduler
	Sequencer         sequence.Sequencer
	Buffer            buffer.Buffer
	Engine            engine.Engine
	Dispatcher        *dispatcher.Dispatcher
	TimeSource        monotime.TimeSource
	statusMutex       sync.RWMutex
	status            PipelineStatus
	Sigma             float64
	AnamolyProbablity float64
	Magnitude         float64
	DriftRate         float64
}
type PipelineStatus struct {
	IsRunning       bool
	EventsProcessed int64
	LastError       error
	StartTime       time.Time
}

type PipelineConfig struct {
	Frequency         event.Frequency
	BufferSize        int
	DLQDirectory      string
	InstanceID        string
	ProducerVersion   string
	TimeSource        monotime.TimeSource
	Dispatcher        dispatcher.DispatcherConfig
	Users             *user.UserRegistry
	Sigma             float64
	AnamolyProbablity float64
	Magnitude         float64
	DriftRate         float64
}

// How FrequencyPipeline will use Transport
func New(cfg PipelineConfig, tsp transport.Transport) (*FrequencyPipeline, error) {
	buf := buffer.New(cfg.BufferSize)
	seq := sequence.New()
	sch := scheduler.New(cfg.Frequency, cfg.TimeSource, cfg.BufferSize)

	d, err := dlq.NewFileDlq(cfg.DLQDirectory, cfg.InstanceID, cfg.TimeSource)
	if err != nil {
		return nil, err
	}

	eng := engine.New(sch, seq, buf, cfg.Users, cfg.ProducerVersion, cfg.InstanceID, cfg.Sigma, cfg.AnamolyProbablity, cfg.Magnitude, cfg.DriftRate)
	ds := dispatcher.New(buf, tsp, cfg.Dispatcher, cfg.TimeSource, d)

	return &FrequencyPipeline{
		Freq:              cfg.Frequency,
		Scheduler:         sch,
		Sequencer:         seq,
		Buffer:            buf,
		Engine:            *eng,
		Dispatcher:        ds,
		TimeSource:        cfg.TimeSource,
		Sigma:             cfg.Sigma,
		AnamolyProbablity: cfg.AnamolyProbablity,
		Magnitude:         cfg.Magnitude,
		DriftRate:         cfg.DriftRate,
	}, nil
}

func (fp *FrequencyPipeline) Start(parentCtx context.Context, message string) {
	fp.statusMutex.Lock()
	if fp.status.IsRunning {
		fp.statusMutex.Unlock()
		return
	}

	fp.status.IsRunning = true
	fp.status.StartTime = fp.TimeSource.Now()
	fp.statusMutex.Unlock()

	pipelineCtx, cancel := context.WithCancel(parentCtx)
	fp.cancel = cancel

	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		fp.Engine.Start(pipelineCtx, message)
	}()

	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		fp.Dispatcher.StartBatch(pipelineCtx)
	}()

	//go fq.Engine.Start(ctx, message)///Reminder that you are joke
	//go fq.Dispatcher.StartBatch(ctx)

}

func (fp *FrequencyPipeline) Stop() {
	fp.statusMutex.Lock()
	if !fp.status.IsRunning {
		fp.statusMutex.Unlock()
		return
	}
	fp.statusMutex.Unlock()

	if fp.cancel != nil {
		fp.cancel()
	}

	fp.wg.Wait()
	fp.statusMutex.Lock()
	fp.status.IsRunning = false
	fp.statusMutex.Unlock()
}
