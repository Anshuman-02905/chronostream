package monotime

import (
	"time"
)

type FakeTimeSource struct {
	current time.Time
}

func NewFakeTimeSource(t time.Time) *FakeTimeSource {
	return &FakeTimeSource{current: t}
}

func (f *FakeTimeSource) Now() time.Time {
	return f.current
}

func (f *FakeTimeSource) Advance(d time.Duration) {
	f.current = f.current.Add(d)
}
