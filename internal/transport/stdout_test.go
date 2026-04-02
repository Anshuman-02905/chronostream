package transport

import (
	"context"
	"testing"

	"github.com/Anshuman-02905/chronostream/internal/event"
)

func TestStdoutTransport_Send(t *testing.T) {
	trans := &StdoutTransport{}
	ctx := context.Background()
	ev := event.Event{
		ID:              "dummy",
		Timestamp:       0,
		Frequency:       event.FrequencySecond,
		Sequence:        0,
		Seed:            0,
		SchemaVersion:   0,
		ProducerVersion: "dummy",
		InstanceID:      "dummy",
		Payload:         []byte("Test-This-Out"),
	}
	err := trans.Send(ctx, ev)
	if err != nil {
		t.Errorf("expected no of error , got %v", err)
	}
}

func TestStdoutTransport_Close(t *testing.T) {
	trans := &StdoutTransport{}
	ctx := context.Background()
	ev := event.Event{
		ID:              "dummy",
		Timestamp:       0,
		Frequency:       event.FrequencySecond,
		Sequence:        0,
		SchemaVersion:   0,
		ProducerVersion: "dummy",
		InstanceID:      "dummy",
		Payload:         []byte("test-this out"),
	}

	_ = trans.Send(ctx, ev)

	err := trans.Close(ctx)
	if err != nil {
		t.Errorf("Expected no eroor but received %v", err)
	}
}
