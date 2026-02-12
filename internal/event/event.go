package event

// Frequ ency Represents how often a event is emitted
// It is intentionally typed

type Frequency uint8

const (
	FrequencyUnknown Frequency = iota
	FrequencySecond
	FrequencyMinute
	FrequencyHour
	FrequencyDay
)
