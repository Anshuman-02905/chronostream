package buffer

import (
	"sync"

	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/sirupsen/logrus"
)

// defined Buffer Interface
// It currenctly has no Transport coupling
// Data Sematics are Pass by Value
type Buffer interface {
	Offer(event.Event) bool
	Events() <-chan event.Event
	Len() int
	Cap() int
	Close()
}

type RealBuffer struct {
	ch   chan event.Event
	once sync.Once
}

// It is a bounded Buffer
func New(capacity int) *RealBuffer {
	logrus.Infof("Creating buffer %v", capacity)

	return &RealBuffer{
		ch: make(chan event.Event, capacity),
	}
}

// It has a non blocking offer
func (r *RealBuffer) Offer(e event.Event) bool {
	logrus.Infof("Buffering")
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

// Len and cap functions are exposing depth introspection
func (r *RealBuffer) Len() int {
	return len(r.ch)
}

func (r *RealBuffer) Cap() int {
	return cap(r.ch)
}

func (r *RealBuffer) Close() {
	r.once.Do(func() {
		close(r.ch)
	})
}
