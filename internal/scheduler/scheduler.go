package scheduler

import (
	"context"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
	"github.com/sirupsen/logrus"
)

//Scheduler Charectersticks
// Depends on monotime.TimeSource interface
// Has a clean contructor
// Computes aligned wall clock boundaries
// Avoids Drift
// Emits deterministic SchduledTime
// Never block on slow consumer
// Respect Context Cancellation
// Has no buffering coupling
// Has no Downstream knowledge

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
	logrus.Infof("Creating Scheduler %v,%v", freq, ts)
	return &RealScheduler{
		frequency: freq,
		ts:        ts,
		ticks:     make(chan Tick, bufferSize),
	}
}

// This helps keep Start()clean
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
		return truncated.Add(durationFor(freq))

	case event.FrequencyMinute:
		truncated := now.Truncate(time.Minute)
		return truncated.Add(durationFor(freq))
	case event.FrequencyHour:
		truncated := now.Truncate(time.Hour)
		return truncated.Add(durationFor(freq))
	case event.FrequencyDay:
		//For day,truncate to midnight in current location
		//We cannot use Truncate(24*time.hour) because 24th from the epoch is not allight to local midnight
		// So we explicitly contruct the midnight in the current location
		//this preserves the timezone DST behaviour Human Expected calendar boundary
		year, month, day := now.Date()
		loc := now.Location()
		midnight := time.Date(year, month, day, 0, 0, 0, 0, loc)
		return midnight.Add(durationFor(freq))
	default:
		panic("unsupported frequency")
	}
}

//design Gaurantees till now
//Process starts at random time
//system clock jitter exists
//Gouroutine wake late

// Start
// DRIFT prevention as if GC pauses System sleeps CPU is overloaded We snap back to correct boundary no cumulative Drift is observed
// Here we spawn  a goroutine
//
//	Loop Forever we check the current time
//	compute the next boundary
//	the sleep duration
//	We wait for the Duration then Emit Tick and Repeat
func (s *RealScheduler) Start(ctx context.Context) {
	go func() {
		for {
			//Get the current Wall clock time
			now := s.ts.Now()

			//Compute the next alligned boundary
			//Always recompute  boundary through recalculation whcih helps us to lock with wall clock forever and avoids DRIFT
			next := nextboundary(now, s.frequency)

			//Calculate how long to wait
			wait := next.Sub(now)
			if wait < 0 {
				wait = 0
			}

			timer := s.ts.NewTimer(wait)

			select {
			//if a service shutdown then context is cancelled Timer is stopped and Goroutine is exited cleanly
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C():
				tick := Tick{
					Frequency:     s.frequency,
					ScheduledTime: next.UnixNano(), // Unix.Nano is used get the exact number of nanoseconds elapsed from January 1, 1970, 00:00:00 UTC
				}
				//this is a non blocing send if consumer is slow we drop the tick we do not delay time. Time cannot wait for consumers
				select {
				case s.ticks <- tick:
				default:

				}

			}

		}
	}()
}

// receive only channel is returned so that consumers can only read
//
//	Consumers cannot write and
//
// consumers cannt close
func (s *RealScheduler) Ticks() <-chan Tick {
	return s.ticks
}

//NEED TO ADD BACKWARD TIME JUMP Resolution
