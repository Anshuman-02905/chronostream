package signal

import (
	"fmt"
	"math"
)

type SignalType string

const (
	Sine       SignalType = "Sine"
	Cosine     SignalType = "cosine"
	SawTooth   SignalType = "sawtooth"
	RandomWalk SignalType = "random_walk"
	Square     SignalType = "square"
)

func Generate(st SignalType, t int64, amplitude, hz float64) (float64, error) {

	switch st {
	case Sine:
		return amplitude * math.Sin(2*math.Pi*hz*float64(t)), nil
	case Cosine:
		return amplitude * math.Cos(2*math.Pi*hz*float64(t)), nil
	case SawTooth:
		return amplitude * (2 * (float64(t)*hz - math.Floor(float64(t)*hz+0.5))), nil
	case Square:
		return amplitude * math.Sin(math.Sin(2*math.Pi*hz*float64(t))), nil
	case RandomWalk:
		return amplitude * math.Sin(2*math.Pi*hz*float64(t)), nil
	}
	return -1, fmt.Errorf("Signal Type not Found")
}

func GetAllSignals() []SignalType {
	return []SignalType{Sine, Cosine, SawTooth, RandomWalk, Square}
}
