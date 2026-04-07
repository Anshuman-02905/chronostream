package signal

import (
	"fmt"
)

type SignalType string

const (
	Sine       SignalType = "Sine"
	Cosine     SignalType = "cosine"
	SawTooth   SignalType = "sawtooth"
	RandomWalk SignalType = "random_walk"
	Square     SignalType = "square"
)

func GetbaseSignal(st SignalType, amplitude, hz float64) (SignalFunc, error) {

	switch st {
	case Sine:
		return SinSignal(amplitude, hz)
	case Cosine:
		return CosSignal(amplitude, hz)
	case SawTooth:
		return SawToothSignal(amplitude, hz)
	case Square:
		return SquareSignal(amplitude, hz)
	case RandomWalk:
		return RandomWalkSignal(amplitude, hz)
	}
	return nil, fmt.Errorf("Signal Type not Found")
}

func CreateSignalPipeline(
	signalType SignalType,
	amplitude float64,
	hz float64,
	sigma float64,
	sampleSeed float64,
	anamolyProbablity float64,
	anamolyMagnitude float64,
	DriftRate float64,
) (SignalFunc, error) {
	base, err := GetbaseSignal(signalType, amplitude, hz)
	if err != nil {
		return nil, err
	}
	return BuildPipeline(
		base,
		Noise(sigma, sampleSeed),
		Drift(DriftRate),
		Anomaly(anamolyProbablity, anamolyMagnitude, sampleSeed),
	), nil
}

func Generate(signalType SignalType,
	t int64,
	amplitude float64,
	hz float64,
	sigma float64,
	sampleSeed float64,
	anamolyProbablity float64,
	anamolyMagnitude float64,
	driftRate float64,
) (float64, error) {
	pipeline, err := CreateSignalPipeline(signalType, amplitude, hz, sigma, sampleSeed, anamolyProbablity, anamolyMagnitude, driftRate)
	if err != nil {
		return 0, err
	}
	return pipeline(t), err
}

func GetAllSignals() []SignalType {
	return []SignalType{Sine, Cosine, SawTooth, RandomWalk, Square}
}
