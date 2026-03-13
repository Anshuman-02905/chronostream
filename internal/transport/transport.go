package transport

import "github.com/Anshuman-02905/chronostream/internal/event"

//Transport defines the I/O contract for delivering  a single detrministic Event
//Implementation should attempt Delivery and return an error
//This implementation should not handle retries; They should just attempt delivery and retun an error on failure

type Transport interface {
	//Send takes a single event and attempts to deliver it
	// It returns an error if the delivery fails
	Send(e event.Event) error

	//Close safely tears down any iunderlying connections( eg closing Kafka Producer, Http Client)
	Close() error
}
