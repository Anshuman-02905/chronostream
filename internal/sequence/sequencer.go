package sequence

import (
	"sync"

	"github.com/Anshuman-02905/chronostream/internal/event"
)

// this does one thing return the next sequence for a given frequency
// no timestamp awareness
// no builder knowledge(event)
// no scheduler knowledge
type Sequnecer interface {
	Next(freq event.Frequency) uint64
}

type RealSequencer struct {
	counters map[event.Frequency]uint64
	mu       sync.Mutex
}

//Behaviur
// Independent counter per frequency
// Monotoninc per frequency
// never resets
// thread safe

func (s *RealSequencer) Next(freq event.Frequency) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[freq]++
	return s.counters[freq]
}

//No Resets if Reset is allowed then idenntity no longer strictly increasing
