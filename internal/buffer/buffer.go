package buffer

import (
	"github.com/Anshuman-02905/chronostream/internal/event"
)

type Buffer interface {
	Offer(event.Event) bool
	Events() <-chan event.Event
	Len() int
	Cap() int
}

type RealBuffer struct {
	ch chan event.Event
}

func New(capacity int) *RealBuffer {
	return &RealBuffer{
		ch: make(chan event.Event, capacity),
	}
}
func (r *RealBuffer) Offer(e event.Event) bool {
	select {
	case r.ch <- e:
		return true
	default:
		return false
	}
}

func (r *RealBuffer) Events() <-chan event.Event {
	return r.ch
}

func (r *RealBuffer) Len() int {
	return len(r.ch)
}

func (r *RealBuffer) Cap() int {
	return cap(r.ch)
}
