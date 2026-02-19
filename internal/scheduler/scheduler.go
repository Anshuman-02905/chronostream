package scheduler

import (
	"context"

	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
)

type Tick struct {
	//Frequency for which this tick was scehduled
	Frequency event.Frequency
	//ScheduledTime is the exact wall clock this tick represents
	// not the time it observed or processed
	ScheduledTime int64
}

type Scheduler interface {
	Start(ctx context.Context)

	Ticks() <-chan Tick
}

//A concrete Scheduler implementation
//This scheduler is the contruct of interface it has
// owns one Frequency
// uses a TimeSource
// It Computes Next Boundary
// Emit Alligned Ticks
// Respect Context Cancellation
// never Blocks
//It should hold
// freqency
//  timeSourece
//  tickChannel
//  and Buffer size if needed

type RealScheduler struct {
	frequency event.Frequency
	ts        monotime.TimeSource
	ticks     chan Tick
}

func New(freq event.Frequency, ts monotime.TimeSource, bufferSize int) *RealScheduler {
	return &RealScheduler{
		frequency: freq,
		ts:        ts,
		ticks:     make(chan Tick, bufferSize),
	}
}

func durationFor()  {}
func Start()        {}
func nextboundary() {}
