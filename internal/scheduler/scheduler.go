package scheduler

import (
	"context"
	"time"

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

// This helps keep Star()clean
//
//	Avoids use of magic number
//	makes future entention easier
//	and Centralizes frequency behaviour
func durationFor(freq event.Frequency) time.Duration {
	switch freq {
	case event.FrequencySecond:
		return time.Second
	case event.FrequencyMinute:
		return time.Minute
	case event.FrequencyHour:
		return time.Hour
	case event.FrequencyDay:
		return 24 * time.Hour
	default:
		panic("unsupported frequency")
	}
}

// principle instead of adding duration to Now truncate to boundary then add + 1 unit this will help us lock to the real wall clock
func nextboundary(now time.Time, freq event.Frequency) time.Time {
	switch freq {
	case event.FrequencySecond:
		truncated := now.Truncate(time.Second)
		return truncated.Add(time.Second)

	}
}

func Start() {}
