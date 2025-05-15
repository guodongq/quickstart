package ddd

import (
	"context"
	"github.com/guodongq/quickstart/pkg/log"
	"time"
)

type Event interface {
	OccurredAt() time.Time
	EventName() string
}

type BaseEvent struct {
	occurredAt time.Time
	name       string
}

func NewBaseEvent(name string) *BaseEvent {
	return &BaseEvent{
		occurredAt: time.Now().UTC(),
		name:       name,
	}
}

func (e *BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *BaseEvent) EventName() string {
	return e.name
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(handler func(ctx context.Context, event Event) error)
}

type LogPublisher struct {
	logger log.Logger
}

func NewLogPublisher(logger log.Logger) *LogPublisher {
	return &LogPublisher{logger: logger}
}

func (p *LogPublisher) Publish(ctx context.Context, event Event) error {
	p.logger.Infof("Event %s occurred at %s", event.EventName(), event.OccurredAt().String())
	return nil
}

func (p *LogPublisher) Subscribe(handler func(ctx context.Context, event Event) error) {
	p.logger.Infof("Subscribed to event %p", handler)
}
