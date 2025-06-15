package events

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/domain/environment"
	"github.com/guodongq/quickstart/pkg/bus"
)

type EnvironmentEventRegistry struct {
	bus             bus.EventBus
	environmentRepo environment.Repository
}

func NewEnvironmentEventRegistry(
	bus bus.EventBus,
	environmentRepo environment.Repository,
) *EnvironmentEventRegistry {
	return &EnvironmentEventRegistry{
		bus:             bus,
		environmentRepo: environmentRepo,
	}
}

func (e *EnvironmentEventRegistry) Register() {
	e.bus.Subscribe(func(ctx context.Context, event *environment.EnvironmentCreatedEvent) error {
		return nil
	})
}
