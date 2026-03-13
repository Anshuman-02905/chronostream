package transport

import (
	"encoding/json"
	"fmt"

	"github.com/Anshuman-02905/chronostream/internal/event"
)

type StdoutTransport struct{}

func (s *StdoutTransport) Send(e event.Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func (s *StdoutTransport) Close() error {
	return nil
}
