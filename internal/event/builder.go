package event

import (
	"fmt"
)

func Build(freq Frequency, ts int64, seq int64) Event {
	id := buildID(freq, ts, seq)

	return Event{
		ID:        id,
		Frequency: freq,
		Timestamp: ts,
		Sequence:  seq,
	}
}

func buildID(freq Frequency, ts int64, seq int64) string {
	return fmt.Sprintf("%d-%d-%d", freq, ts, seq)
}
