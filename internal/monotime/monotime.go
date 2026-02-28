package monotime

import (
	"time"
)

type Timer interface {
	C() <-chan time.Time
	Stop() bool
}

type TimeSource interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
}

type realTimer struct {
	t *time.Timer
}

func (rt *realTimer) C() <-chan time.Time {
	return rt.t.C
}

func (rt *realTimer) Stop() bool {
	return rt.t.Stop()
}

type RealTimeSource struct{}

func (t *RealTimeSource) Now() time.Time {
	return time.Now().UTC()
}

func (t *RealTimeSource) NewTimer(d time.Duration) Timer {
	return &realTimer{t: time.NewTimer(d)}
}
