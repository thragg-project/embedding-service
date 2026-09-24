package services

import (
	"context"
)

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type Event struct {
	AggregateType string
	AggregateID   int64
	EventType     string
	Payload       any
}
