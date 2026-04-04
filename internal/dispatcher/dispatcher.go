package dispatcher

import (
	"context"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/dlq"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
	"github.com/Anshuman-02905/chronostream/internal/transport"

	"github.com/sirupsen/logrus"
)

// DispatcherConfig holds per-frequency dispatch settings
// Each FrequencyPipeline creates its own DispatcherConfig from the per-frequency config
type DispatcherConfig struct {
	MaxRetries    int
	BaseBackoff   int // milliseconds
	MaxBackoff    int // milliseconds
	BatchSize     int
	FlushInterval int // milliseconds
}

//Dispatcher is responsible for consuming events from buffer and routing them to a transport
//it will eventually handle retries ,Dead letter  and delivery semantics

type Dispatcher struct {
	buf   buffer.Buffer
	trans transport.Transport
	cfg   DispatcherConfig
	ts    monotime.TimeSource
	dlq   dlq.DLQ
}

//New Creates a new Dispatcher wiring a buffer to transport

func New(buf buffer.Buffer, trans transport.Transport, cfg DispatcherConfig, ts monotime.TimeSource, dlq dlq.DLQ) *Dispatcher {
	return &Dispatcher{
		buf:   buf,
		trans: trans,
		cfg:   cfg,
		ts:    ts,
		dlq:   dlq,
	}
}

//Start begins  a blocking loop consuming events from the buffer
//It stops when the context is cancelled or the buffer channel is closed

func (d *Dispatcher) Start(ctx context.Context) {
	baseDelay := time.Duration(d.cfg.BaseBackoff) * time.Millisecond
	maxDelay := time.Duration(d.cfg.MaxBackoff) * time.Millisecond
	maxRetries := d.cfg.MaxRetries
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-d.buf.Events():
			if !ok {
				//buffer closed
				return
			}
			for attempt := 0; attempt <= maxRetries; attempt++ {

				//Attempt to  send  the event via Transport
				err := d.trans.Send(ctx, ev)
				if err == nil {
					break
				}
				if attempt == maxRetries {
					logrus.WithError(err).Error("Max retries  reached , dropping event Redirecting to DLQ")
					d.dlq.Writebatch(ctx, []event.Event{ev})
				}
				backoff := baseDelay * time.Duration(1<<attempt)
				if backoff > maxDelay {
					backoff = maxDelay
				}
				logrus.WithFields(logrus.Fields{
					"attempt": attempt,
					"delay":   backoff,
				}).Warn("Dispatch Field,retrying")

				select {
				case <-time.After(backoff):
				case <-ctx.Done():
					return
				}

			}
		}
	}
}

func (d *Dispatcher) StartBatch(ctx context.Context) {

	flushInterval := time.Duration(d.cfg.FlushInterval) * time.Millisecond
	batchSize := d.cfg.BatchSize
	var batch []event.Event

	timer := d.ts.NewTimer(flushInterval)
	defer timer.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		toSend := batch
		batch = nil
		d.SendBatchWithRetry(ctx, toSend)
	}
	resetTimer := func() {
		timer.Stop()
		timer = d.ts.NewTimer(flushInterval)
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case ev, ok := <-d.buf.Events():
			if !ok {
				flush()
				return
			}
			batch = append(batch, ev)

		drainLoop:
			for len(batch) < batchSize {
				select {
				case nextEv, ok := <-d.buf.Events():
					if !ok {
						flush()
						break drainLoop
					}
					batch = append(batch, nextEv)
				default:
					break drainLoop
				}
			}
			if len(batch) >= batchSize {
				flush()
				resetTimer()
			}

		case <-timer.C():
			flush()
			resetTimer()

		}
	}
}

func (d *Dispatcher) SendBatchWithRetry(ctx context.Context, events []event.Event) {
	baseDelay := time.Duration(d.cfg.BaseBackoff) * time.Millisecond
	maxDelay := time.Duration(d.cfg.MaxBackoff) * time.Millisecond
	maxRetries := d.cfg.MaxRetries

	for attempt := 0; attempt <= maxRetries; attempt++ {
		//Attempt to send event via Transport
		err := d.trans.SendBatch(ctx, events)
		if err == nil {
			break
		}
		if attempt == maxRetries {
			logrus.WithError(err).Error("Max retries reached dropping batch redirecting to DLQ")
			d.dlq.Writebatch(ctx, events)
		}
		backoff := baseDelay * time.Duration(1<<attempt)
		if backoff > maxDelay {
			backoff = maxDelay
		}
		logrus.WithFields(logrus.Fields{
			"attempt": attempt,
			"delay":   backoff,
		}).Warn("Batch dispatch failed Retrying")
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return
		}
	}

}
