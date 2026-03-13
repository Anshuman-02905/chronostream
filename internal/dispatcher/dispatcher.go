package dispatcher

import (
	"context"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/transport"
	"github.com/sirupsen/logrus"
)

//Dispatcher is responsible for consuming events from buffer and routing them to a transport
//it will eventually handle retries ,Dead letter  and delivery semantics

type Dispatcher struct {
	buf   buffer.Buffer
	trans transport.Transport
}

//New Creates a new Dispatcher wiring a buffer to transport

func New(buf buffer.Buffer, trans transport.Transport) *Dispatcher {
	return &Dispatcher{
		buf:   buf,
		trans: trans,
	}
}

//Start begins  a blocking loop consuming events from the buffer
//It stops when the context is cancelled or the buffer channel is closed

func (d *Dispatcher) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-d.buf.Events():
			if !ok {
				//buffer closed
				return
			}

			//Attempt to  send  the event via Transport
			err := d.trans.Send(ev)
			if err != nil {
				//For milestone 1 , just log the error
				//Later we implementation exponential backoff / retries here
				logrus.Errorf("Failed to Dispatch the event")
			}

		}
	}
}
