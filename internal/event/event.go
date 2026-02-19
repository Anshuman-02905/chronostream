package event

// Frequency Represents how often a event is emitted
// It is intentionally typed (not string) to guarantee
//compile-safety
// Exhastive Switch handling
// determinstic behaviour

type Frequency uint8

const (
	FrequencyUnknown Frequency = iota
	FrequencySecond
	FrequencyMinute
	FrequencyHour
	FrequencyDay
)

// Event is the immutable core unit of the system
//
//Design gurantees:
// - Determinitic foreever give (timestamp, frequency , sequence)
// - No pointers
// - no runtime depencency (time.now, randomness, IO)
// No transport or Serization concerns
//
//Once created this struct must never change the shape lightly

type Event struct {
	//Deterministeic Identity
	ID string

	//Core Timing
	////unix timestamp in nanoseconds (explicit and unambigous)
	Timestamp int64
	////Frequency at which this event was guaranteed
	Frequency Frequency
	////Sequence number within the same timestamp+frequency window
	Sequence uint64

	//Deterministic Seed for downstream systems
	Seed int64
	//Schema Evolution Safety
	SchemaVersion uint16
	//Producer Metadata (flat ,immutable)
	ProducerVersion string
	InstanceID      string
}
