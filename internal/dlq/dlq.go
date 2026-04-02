package dlq

import (
	"context"

	"github.com/Anshuman-02905/chronostream/internal/event"
)

type DLQ interface {
	//WriteBatch appends a failed batch of events to a fallback storage
	Writebatch(ctx context.Context, events []event.Event) error
	//RouteToDLQ  Routes a batch toward DLQ Write Batch Checks if dlq is initialized or not then only it redirects
	Close(ctx context.Context) error
}
