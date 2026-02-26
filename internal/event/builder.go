package event

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Build Contructs a fully deterministic , immutatble Event

// Determinism Contract:
// Givent the same (Frequency Timestamp Sequence producerVersion InstanceID)
// this function MUST always return an indentical Event

// Design Rules:
// 	-No Randomness Used
// 	-No time.Now
// 	-No Global State
// 	-No I/O operation
// 	-No Global State
// 	-Pure function

// The Builder is the single authority responsible for :
// 	- Event identity Contstruction
// 	- Deterministic seed derivation
// 	- Schema Version Stamping

// Any Change to ID or Seed generation Logic is a breakig change for downstream consumers

func Build(freq Frequency, ts int64, seq uint64, producerVersion string, instanceID string) Event {
	id := buildID(freq, ts, seq)
	seed := buildSeed(ts, seq)
	logrus.Infof("Building event %v,%v", freq, ts)

	return Event{
		ID:              id,
		Frequency:       freq,
		Timestamp:       ts,
		Sequence:        seq,
		Seed:            seed,
		SchemaVersion:   1,
		ProducerVersion: producerVersion,
		InstanceID:      instanceID,
	}
}

func buildID(freq Frequency, ts int64, seq uint64) string {
	return fmt.Sprintf("%d-%d-%d", freq, ts, seq)
}

func buildSeed(ts int64, seq uint64) int64 {
	return ts ^ int64(seq)
}
