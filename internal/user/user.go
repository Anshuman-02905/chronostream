package user

import (
	"fmt"

	"github.com/Anshuman-02905/chronostream/internal/signal"
)

type User struct {
	ID         string            // "user_001" "user_002", etc
	Session    string            // Unique string
	SignalType signal.SignalType // "sine" , "cosine", "sawtooth" etc
}

func NewUser(id string, session string, st signal.SignalType) (*User, error) {
	// IsValidSignalType returns true for valid types, so we error on !valid
	if !IsValidSignalType(st) {
		return nil, fmt.Errorf("incorrect signal type: %q", st)
	}
	return &User{
		ID:         id,
		Session:    session,
		SignalType: st,
	}, nil
}

// IsValidSignalType returns true if st is a known signal type
func IsValidSignalType(st signal.SignalType) bool {
	for _, s := range signal.GetAllSignals() {
		if st == s {
			return true
		}
	}
	return false
}
