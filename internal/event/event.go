package event

import "fmt"

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

type EventType uint8

const (
	EventTypeUnknown          EventType = iota
	EventTypeChatMessage                // FrequencySecond
	EventTypeAggregatedMetric           // FrequencyMinute
	EventTypeSessionSnapshot            // FrequencyHour
	EventTypeDailyMarker                // FrequencyDay
)

func EventTypeFor(freq Frequency) EventType {
	switch freq {
	case FrequencySecond:
		return EventTypeChatMessage
	case FrequencyMinute:
		return EventTypeAggregatedMetric
	case FrequencyHour:
		return EventTypeSessionSnapshot
	case FrequencyDay:
		return EventTypeUnknown
	default:
		return EventTypeUnknown
	}
}

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
	EventType EventType

	ProducerVersion string
	InstanceID      string
	Payload         []byte

	UserID    string
	SessionID string

	MessageID     string
	FragmentIndex int

	TotalFragments int
}

// ParseFrequency converts a config string ("second", "minute", "hour", "day")
// to the typed Frequency enum. Returns error for unknown strings.
func ParseFrequency(s string) (Frequency, error) {
	switch s {
	case "second":
		return FrequencySecond, nil
	case "minute":
		return FrequencyMinute, nil
	case "hour":
		return FrequencyHour, nil
	case "day":
		return FrequencyDay, nil
	default:
		return FrequencyUnknown, fmt.Errorf("unknown frequency: %q", s)
	}
}
