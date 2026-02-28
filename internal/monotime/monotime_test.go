package monotime

import (
	"testing"
	"time"
)

func TestRealTimeSource_Now_ReturnsUTC(t *testing.T) {
	ts := &RealTimeSource{}
	now := ts.Now()

	if now.Location() != time.UTC {
		t.Fatalf("Expected UTC location but got %v ", now.Location())
	}
}

func TestRealTimeSource_Now_IsCloseToSystemTime(t *testing.T) {
	ts := &RealTimeSource{}
	before := time.Now().UTC()
	now := ts.Now()
	after := time.Now().UTC()

	if now.Before(before) || now.After(after) {
		t.Fatalf("Now () returned outside expected Range")
	}
}

// Think the Abstaction should work correctly
func TestRealTimer(t *testing.T) {

	ts := &RealTimeSource{}
	d := time.Duration(2 * time.Second)
	tm := ts.NewTimer(d)

	start := time.Now()

	select {
	case <-tm.C():
		elapsed := time.Since(start)
		if elapsed < 2*time.Second {
			t.Fatalf("timer fired too early: %v", elapsed)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timer did not fire within expected time")

	}
}

func TestRealTimer_Stop_BeforeFire(t *testing.T) {
	ts := &RealTimeSource{}
	d := time.Duration(2 * time.Second)

	tm := ts.NewTimer(d)
	stopped := tm.Stop()
	if !stopped {
		t.Fatalf("Expected Stop() to return true")
	}

	select {
	case <-tm.C():
		t.Fatalf("Timer fired after Stop()")
	case <-time.After(3 * time.Second):
	}
}
