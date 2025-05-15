package event

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/domain/environment"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/bus"
)

type EventRegistry interface {
	RegisterHandlers()
}

type ProjectEventRegistry struct {
	bus         bus.EventBus
	projectRepo project.Repository
}

func NewProjectEventRegistry(
	bus bus.EventBus,
	projectRepo project.Repository,
) *ProjectEventRegistry {
	return &ProjectEventRegistry{
		bus:         bus,
		projectRepo: projectRepo,
	}
}

func (e *ProjectEventRegistry) RegisterHandlers() {
	e.bus.Subscribe(func(ctx context.Context, event *project.ProjectCreatedEvent) error {
		return nil
	})
}

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

func (e *EnvironmentEventRegistry) RegisterHandlers() {
	e.bus.Subscribe(func(ctx context.Context, event *environment.EnvironmentCreatedEvent) error {
		return nil
	})
}
