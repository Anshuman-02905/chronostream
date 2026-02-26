package sequence

import (
	"testing"

	"github.com/Anshuman-02905/chronostream/internal/event"
)

func TestNext(t *testing.T) {

	ts := New()
	evf := event.FrequencySecond
	seq := ts.Next(evf)
	expected_seq := uint64(1)

	if seq != expected_seq {
		t.Fatalf("Expected seq %v got %v", expected_seq, seq)
	}

}
