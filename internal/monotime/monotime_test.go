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
