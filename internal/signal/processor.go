package signal

import (
	"math"
	"math/rand/v2"
)

type SignalFunc func(t int64) float64
type Decorator func(SignalFunc) SignalFunc

func BuildPipeline(
	base SignalFunc,
	decorators ...Decorator,
) SignalFunc {
	s := base
	for _, d := range decorators {
		s = d(s)
	}
	return s
}

func SinSignal(amplitude, hz float64) (SignalFunc, error) {
	return func(t int64) float64 {
		return amplitude * math.Sin(2*math.Pi*hz*float64(t))
	}, nil
}

func CosSignal(amplitude, hz float64) (SignalFunc, error) {
	return func(t int64) float64 {
		return amplitude * math.Sin(2*math.Pi*hz*float64(t))
	}, nil
}

func SawToothSignal(amplitude, hz float64) (SignalFunc, error) {
	return func(t int64) float64 {
		return amplitude * (2 * (float64(t)*hz - math.Floor(float64(t)*hz+0.5)))
	}, nil
}

func SquareSignal(amplitude, hz float64) (SignalFunc, error) {
	return func(t int64) float64 {
		return amplitude * math.Sin(math.Sin(2*math.Pi*hz*float64(t)))
	}, nil
}
func RandomWalkSignal(amplitude, hz float64) (SignalFunc, error) {
	return func(t int64) float64 {
		return amplitude * math.Sin(2*math.Pi*hz*float64(t))
	}, nil
}

func Drift(driftRate float64) Decorator {
	return func(base SignalFunc) SignalFunc {
		return func(t int64) float64 {
			drift := driftRate * float64(t)
			return base(t) + drift
		}
	}
}

func Anomaly(probablity float64, magnitude float64, seed float64) Decorator {
	return func(base SignalFunc) SignalFunc {
		rng := rand.New(rand.NewPCG(uint64(seed), uint64(seed+1)))
		return func(t int64) float64 {
			value := base(t)
			//Roll Dice
			if rng.Float64() < probablity {
				spike := magnitude
				if rand.IntN(2) == 0 {
					spike = -magnitude
				}
				value += spike
			}
			return value
		}
	}
}

func Noise(sigma float64, seed float64) Decorator {
	return func(base SignalFunc) SignalFunc {
		rng := rand.New(rand.NewPCG(uint64(seed), uint64(seed+1)))
		return func(t int64) float64 {
			value := base(t)
			noise := GausianRandom(rng, sigma)
			return value + noise
		}
	}

}
func GausianRandom(rng *rand.Rand, stdDev float64, mean ...float64) float64 {
	m := 0.0
	if len(mean) > 0 {
		m = mean[0]
	}
	noise := m + stdDev*rng.NormFloat64()
	return noise
}
