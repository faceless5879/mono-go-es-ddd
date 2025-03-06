package event_sourcing

import "github.com/google/uuid"

type AggregateRoot interface {
	GetUncommitedEvents() []Event
	ID() uuid.UUID
	LoadFromHistory(events []Event)
	ApplyChange(event Event, loadFromHistory bool)
	Version() int
}
