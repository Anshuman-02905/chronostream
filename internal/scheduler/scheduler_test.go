package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
)

func TestNextBoundarySecond(t *testing.T) {
	now := time.Date(2026, 2, 20, 10, 15, 42, 800_000_000, time.UTC)

	next := nextboundary(now, event.FrequencySecond)

	expected := time.Date(2026, 2, 20, 10, 15, 43, 0, time.UTC)

	if !next.Equal(expected) {
		t.Fatalf("expected %v got %v", expected, next)
	}
}

func TestNextBoundaryMinute(t *testing.T) {
	now := time.Date(2026, 2, 20, 10, 15, 42, 800_000_000, time.UTC)
	next := nextboundary(now, event.FrequencyMinute)
	expected := time.Date(2026, 2, 20, 10, 16, 00, 0, time.UTC)

	if !next.Equal(expected) {
		t.Fatalf("expected %v got %v", expected, next)

	}
}

func TestNextBoundaryHour(t *testing.T) {
	now := time.Date(2026, 2, 20, 10, 15, 42, 800_000_000, time.UTC)
	next := nextboundary(now, event.FrequencyHour)
	expected := time.Date(2026, 2, 20, 11, 00, 0, 0, time.UTC)

	if !next.Equal(expected) {
		t.Fatalf("expected %v got %v ", expected, next)
	}
}

func TestNextBoundaryDay(t *testing.T) {
	now := time.Date(2026, 2, 20, 10, 15, 42, 800_000_000, time.UTC)
	next := nextboundary(now, event.FrequencyDay)
	expected := time.Date(2026, 2, 21, 00, 00, 00, 0, time.UTC)

	if !next.Equal(expected) {
		t.Fatalf("expected %v got %v ", expected, next)
	}
}

func TestStart(t *testing.T) {
	start := time.Date(2026, 2, 20, 10, 15, 42, 800_000_000, time.UTC)
	fake := monotime.NewFakeTimeSource(start)
	s := New(event.FrequencySecond, fake, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	tick := <-s.Ticks()
	expected := time.Date(2026, 2, 20, 10, 15, 43, 0, time.UTC)

	if tick.ScheduledTime != expected.UnixNano() {
		t.Fatalf("expected %v got %v", expected, time.Unix(0, tick.ScheduledTime))
	}
}
func TestStart_NoConsumer(t *testing.T) {
	done := make(chan struct{})
	go func() {
		start := time.Date(2026, 2, 20, 10, 15, 42, 8000_000_000, time.UTC)
		fake := monotime.NewFakeTimeSource(start)
		s := New(event.FrequencySecond, fake, 1)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		s.Start(ctx)
		for i := 0; i < 5; i++ {
			fake.Advance(time.Second)
		}
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("")
	case <-time.After(1 * time.Second):
		t.Fatalf("Scheduler deadlocked when channel not consumed")
	}

}
