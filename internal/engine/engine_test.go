package engine

import (
	"context"
	"testing"
	"time"

	"github.com/Anshuman-02905/chronostream/internal/buffer"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
	"github.com/Anshuman-02905/chronostream/internal/scheduler"
	"github.com/Anshuman-02905/chronostream/internal/sequence"
)

func TestEngine_EndToEnd(t *testing.T) {
	start := time.Now().Truncate(time.Second)
	fakeTime := monotime.NewFakeTimeSource(start)

	sch := scheduler.New(event.FrequencySecond, fakeTime, 1)
	seq := sequence.New()
	buf := buffer.New(10)
	prod_version := "v1.0"
	instance_id := "02905"
	e := New(sch, seq, buf, prod_version, instance_id)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e.Start(ctx)

	fakeTime.Advance(time.Second)

	select {
	case ev := <-buf.Events():
		if ev.Sequence != 1 {
			t.Errorf("Expected sequence 1 but got %v", ev.Sequence)
		}
		if ev.Frequency != event.FrequencySecond {
			t.Errorf("Wrong Frequency")
		}
	case <-time.After(time.Millisecond * 500):
		t.Fatalf("Timeout:Event Never reached the buffer")
	}

}
