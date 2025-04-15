package order

import (
	"context"
	"fmt"
	es "github.com/faceless5879/mono-go-es-ddd/internal/common/event_sourcing"
)

type NotFoundError struct {
	OrderID string
}

func (e NotFoundError) Error() string { return fmt.Sprintf("training '%s' not found", e.OrderID) }

type Repository interface {
	SaveEvents(ctx context.Context, events []es.Event) error
	LoadEvents(ctx context.Context, aggregateID string) ([]es.Event, error)
}
