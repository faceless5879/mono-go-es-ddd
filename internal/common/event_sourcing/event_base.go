package event_sourcing

import (
	"time"
)

type Event interface {
	EventType() EventType
	TimeStamp() time.Time
}

type BaseEvent struct {
	EventType EventType
	TimeStamp time.Time
}

type EventType string
