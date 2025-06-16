package commands

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/bus"
	"github.com/guodongq/quickstart/pkg/ddd"
	"github.com/guodongq/quickstart/pkg/decorator"
	"github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/types"
)

type CreateProjectHandler decorator.CommandHandler[dto.CreateProject]

type createProjectHandler struct {
	projectRepository project.Repository
	eventBus          bus.EventBus
}

func NewCreateProjectHandler(
	projectRepository project.Repository,
	logger log.Logger,
	metricsClient decorator.MetricsClient,
	eventBus bus.EventBus,
) CreateProjectHandler {
	return decorator.ApplyCommandDecorators[dto.CreateProject](
		createProjectHandler{
			projectRepository: projectRepository,
			eventBus:          eventBus,
		},
		logger,
		metricsClient,
	)
}

func (c createProjectHandler) Handle(ctx context.Context, cmd dto.CreateProject) error {
	entity := project.NewProject(
		project.WithProjectBaseEntity(ddd.NewBaseEntity(cmd.ID)),
		project.WithProjectName(cmd.Name),
		project.WithProjectDescription(cmd.Description),
		project.WithProjectMeta(types.NewMetaWithDefaults(ctx)),
	)

	if err := c.projectRepository.CreateProject(ctx, entity); err != nil {
		return err
	}

	if err := c.eventBus.Publish(ctx, project.NewProjectCreatedEvent(entity.Id, entity.Name)); err != nil {
		return err
	}

	return nil
}
