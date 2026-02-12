package scheduler

import (
	"context"

	"github.com/Anshuman-02905/chronostream/internal/event"
)

type Tick struct {
	//Frequency for which this tick was scehduled
	Frqquency event.Frequency
	//ScheduledTime is the exact wall clock this tick represents
	// not theß time it observed or processed
	ScheduledTime int64
}

type Scheduler interface {
	Start(ctx context.Context)

	Ticks() <-chan Tick
}
