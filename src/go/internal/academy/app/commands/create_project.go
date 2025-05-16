package commands

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/decorator"
	"github.com/guodongq/quickstart/pkg/log"
)

type CreateProjectHandler decorator.CommandHandler[dto.CreateProject]

type createProjectHandler struct {
	projectRepository project.Repository
}

func NewCreateProjectHandler(
	projectRepository project.Repository,
	logger log.Logger,
	metricsClient decorator.MetricsClient,
) CreateProjectHandler {
	return decorator.ApplyCommandDecorators[dto.CreateProject](
		createProjectHandler{
			projectRepository: projectRepository,
		},
		logger,
		metricsClient,
	)
}

func (c createProjectHandler) Handle(ctx context.Context, cmd dto.CreateProject) error {
	//TODO implement me
	panic("implement me")
}
