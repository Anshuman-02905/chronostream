package engine

import (
	"context"
	"encoding/json"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/chunker"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/scheduler"
	"github.com/Anshuman-02905/chronostream/internal/sequence"
	"github.com/Anshuman-02905/chronostream/internal/signal"
	"github.com/Anshuman-02905/chronostream/internal/user"
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
	registry  *user.UserRegistry

	producerVersion string
	instanceID      string
}

type UserSignalPayload struct {
	UserID    string  `json:"user_id"`
	Session   string  `json:"session"`
	Signal    string  `json:"signal"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

func New(
	s scheduler.Scheduler,
	seq sequence.Sequencer,
	buf buffer.Buffer,
	registry *user.UserRegistry,
	producerVersion string,
	instanceID string,
) *Engine {

	return &Engine{
		scheduler:       s,
		sequencer:       seq,
		buffer:          buf,
		registry:        registry,
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
func (e *Engine) Start(ctx context.Context, message string) {
	users := e.registry.All()
	logrus.WithField("user_count", len(users)).Info("Engine starting")

	// Guard: no users means no events will ever be emitted
	if len(users) == 0 {
		logrus.Error("Engine has 0 users in registry — check users.count in config.yaml")
		return
	}

	e.scheduler.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				e.buffer.Close()
				return
			case tick, ok := <-e.scheduler.Ticks():
				if !ok {
					logrus.Warn("Scheduler ticks channel closed")
					return
				}

				logrus.WithFields(logrus.Fields{
					"frequency":  tick.Frequency,
					"user_count": len(users),
				}).Debug("Tick received, emitting events for all users")

				// Emit one event per user per tick
				for _, u := range users {
					tSec := float64(tick.ScheduledTime) / 1e9
					yValue, err := signal.Generate(u.SignalType, int64(tSec), 1.0, 1.0)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"user_id":     u.ID,
							"signal_type": u.SignalType,
						}).WithError(err).Error("Signal generation failed — skipping user this tick")
						continue
					}

					p := UserSignalPayload{
						UserID:    u.ID,
						Session:   u.Session,
						Signal:    string(u.SignalType),
						Value:     yValue,
						Timestamp: int64(tSec),
					}
					jsonBytes, err := json.Marshal(p)
					if err != nil {
						logrus.WithField("user_id", u.ID).WithError(err).Error("JSON marshal failed — skipping user this tick")
						continue
					}

					fragments := chunker.Chunk(string(jsonBytes), 1024)
					for _, frag := range fragments {
						seq := e.sequencer.Next(tick.Frequency)
						ev := event.Build(
							tick.Frequency,
							tick.ScheduledTime,
							seq,
							e.producerVersion,
							e.instanceID,
							frag.Payload,
							u.Session,
							u.ID,
							frag.MessageID,
							frag.ChunkIndex,
							frag.TotalChunks,
						)
						e.buffer.Offer(ev)
					}

					logrus.WithFields(logrus.Fields{
						"user_id":     u.ID,
						"signal_type": u.SignalType,
						"value":       yValue,
						"frequency":   tick.Frequency,
					}).Debug("Event emitted")
				}
			}
		}
	}()
}
