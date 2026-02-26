package engine

import (
	"context"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/scheduler"
	"github.com/Anshuman-02905/chronostream/internal/sequence"
	"github.com/sirupsen/logrus"
)

//Engine is the composition root
// A wiring layer
// A lifecycle controller // Coordinator
//It consumes the interfaces it does not need to expose an interface
//having an interface engine is wasteful abstraction

type Engine struct {
	scheduler scheduler.Scheduler
	sequencer sequence.Sequencer
	buffer    buffer.Buffer

	producerVersion string
	instanceID      string
}

func New(
	s scheduler.Scheduler,
	seq sequence.Sequencer,
	buf buffer.Buffer,
	producerVersion string,
	instanceID string,
) *Engine {
	logrus.Infof("Creating Engine event %v,%v,%v,%v,%v", s, seq, buf, producerVersion, instanceID)

	return &Engine{
		scheduler:       s,
		sequencer:       seq,
		buffer:          buf,
		producerVersion: producerVersion,
		instanceID:      instanceID,
	}
}

// Engine owns Dependencies MetaData Constants and Lifecycle
// Engine does not compute Boundaries
// Does not compute Sequences
// Does not Mutate events
// Does not block Scheduler
// Does not Know Transport
// Only Wires/orchaestrate/glue services
func (e *Engine) Start(ctx context.Context) {
	e.scheduler.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tick, ok := <-e.scheduler.Ticks():
				if !ok {
					return
				}
				seq := e.sequencer.Next(tick.Frequency)
				ev := event.Build(
					tick.Frequency,
					tick.ScheduledTime,
					seq,
					e.producerVersion,
					e.instanceID,
				)
				e.buffer.Offer(ev)
			}
		}
	}()

}
