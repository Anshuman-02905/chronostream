package dispatcher

import (
	"context"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/config"
	"github.com/Anshuman-02905/chronostream/internal/transport"
	"github.com/sirupsen/logrus"
)

//Dispatcher is responsible for consuming events from buffer and routing them to a transport
//it will eventually handle retries ,Dead letter  and delivery semantics

type Dispatcher struct {
	buf   buffer.Buffer
	trans transport.Transport
	cfg   config.Config
}

//New Creates a new Dispatcher wiring a buffer to transport

func New(buf buffer.Buffer, trans transport.Transport, cfg config.Config) *Dispatcher {
	return &Dispatcher{
		buf:   buf,
		trans: trans,
		cfg:   cfg,
	}
}

//Start begins  a blocking loop consuming events from the buffer
//It stops when the context is cancelled or the buffer channel is closed

func (d *Dispatcher) Start(ctx context.Context) {
	baseDelay := time.Duration(d.cfg.Dispatcher.BaseBackoff) * time.Millisecond
	maxDelay := time.Duration(d.cfg.Dispatcher.MaxBackoff) * time.Millisecond
	maxRetries := d.cfg.Dispatcher.MaxRetries
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
					logrus.WithError(err).Error("Max retries  reached , dropping event")
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
