package event

import (
	"fmt"
)

func Build(freq Frequency, ts int64, seq uint64) Event {
	id := buildID(freq, ts, seq)
	seed := buildSeed(ts, seq)
	return Event{
		ID:        id,
		Frequency: freq,
		Timestamp: ts,
		Sequence:  seq,
		Seed:      seed,
	}
}

func buildID(freq Frequency, ts int64, seq uint64) string {
	return fmt.Sprintf("%d-%d-%d", freq, ts, seq)
}

func buildSeed(ts int64, seq uint64) int64 {
	return ts ^ int64(seq)
}
