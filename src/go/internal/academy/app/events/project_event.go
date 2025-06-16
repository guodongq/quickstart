package events

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/bus"
	logger "github.com/guodongq/quickstart/pkg/log"
)

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

func (e *ProjectEventRegistry) Register() {
	e.bus.Subscribe(func(ctx context.Context, event *project.ProjectCreatedEvent) error {
		logger.Info("Received project.created event: %s", event)
		return nil
	})
}
