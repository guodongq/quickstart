package commands

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/decorator"
	"github.com/guodongq/quickstart/pkg/log"
)

type DeleteProjectHandler decorator.CommandHandler[dto.DeleteProject]

type deleteProjectHandler struct {
	projectRepository project.Repository
}

func NewDeleteProjectHandler(
	projectRepository project.Repository,
	logger log.Logger,
	metricsClient decorator.MetricsClient,
) DeleteProjectHandler {
	return decorator.ApplyCommandDecorators[dto.DeleteProject](
		deleteProjectHandler{
			projectRepository: projectRepository,
		},
		logger,
		metricsClient,
	)
}

func (c deleteProjectHandler) Handle(ctx context.Context, cmd dto.DeleteProject) error {
	//TODO implement me
	panic("implement me")
}
