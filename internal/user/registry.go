package user

import (
	"fmt"
	"math/rand"

	"github.com/Anshuman-02905/chronostream/internal/signal"
)

type UserRegistry struct {
	users map[string]*User
	count int
	seed  int64
}

func NewUserRegistry(count int, seed int64) (*UserRegistry, error) {

	ur := &UserRegistry{
		users: make(map[string]*User),
		count: count,
		seed:  seed,
	}
	// Deterministic random Source
	r := rand.New(rand.NewSource(seed))

	//Available signal types
	signals := signal.GetAllSignals()

	//Create users
	for i := 1; i <= count; i++ {
		id := fmt.Sprintf("user_%03d", i)

		//Determistic Session
		session := fmt.Sprintf("sesssion_%06d", r.Intn(1000000))

		//deteministic Signal
		signal := signals[r.Intn(len(signals))]

		user, err := NewUser(
			id,
			session,
			signal,
		)
		if err != nil {
			return nil, fmt.Errorf("The User was not created")
		}
		ur.users[id] = user
	}
	return ur, nil
}

func (ur *UserRegistry) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("User Id cannot be empty\n")
	}
	user, ok := ur.users[id]
	if !ok {
		return nil, fmt.Errorf("User  not found \n")
	}
	return user, nil
}
func (ur *UserRegistry) AssignSignal(id string, signalType signal.SignalType) error {
	if id == "" {
		return fmt.Errorf("User Id was nil\n")
	}
	if signalType == "" {
		return fmt.Errorf("Signal type was nil")
	}
	//Create a User
	user, ok := ur.users[id]
	if !ok {
		return fmt.Errorf("User was not fould please create")
	}
	//Assign the signal type
	user.SignalType = signalType
	return nil

}

func (ur *UserRegistry) All() []*User {
	result := make([]*User, 0, len(ur.users))
	for _, u := range ur.users {
		result = append(result, u)
	}
	return result
}
