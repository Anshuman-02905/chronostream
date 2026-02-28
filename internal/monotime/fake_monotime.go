package monotime

import (
	"time"
)

type fakeTimer struct {
	c      chan time.Time
	fireAt time.Time
	active bool
}

func (ft *fakeTimer) C() <-chan time.Time {
	return ft.c
}

func (ft *fakeTimer) Stop() bool {
	wasActive := ft.active
	ft.active = false
	return wasActive
}

type FakeTimeSource struct {
	current time.Time
	timers  []*fakeTimer
}

func NewFakeTimeSource(t time.Time) *FakeTimeSource {
	return &FakeTimeSource{
		current: t,
		timers:  make([]*fakeTimer, 0),
	}
}

func (f *FakeTimeSource) Now() time.Time {
	return f.current
}

func (f *FakeTimeSource) NewTimer(d time.Duration) Timer {
	fireAt := f.current.Add(d)
	ft := &fakeTimer{
		c:      make(chan time.Time, 1),
		fireAt: fireAt,
		active: true,
	}

	// If the requested time has already passed (due to Advance being called
	// between the scheduler's Now() call and this NewTimer() call), fire it immediately.
	if !ft.fireAt.After(f.current) {
		ft.c <- ft.fireAt
		ft.active = false
	}

	f.timers = append(f.timers, ft)
	return ft
}

func (f *FakeTimeSource) Advance(d time.Duration) {
	newTime := f.current.Add(d)

	// Check for timers that should fire between current and newTime
	for _, ft := range f.timers {
		if ft.active && !ft.fireAt.After(newTime) {
			ft.c <- ft.fireAt
			ft.active = false
		}
	}

	f.current = newTime
}
