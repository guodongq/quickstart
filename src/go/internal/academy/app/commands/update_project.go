package commands

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/decorator"
	"github.com/guodongq/quickstart/pkg/log"
)

type UpdateProjectHandler decorator.CommandHandler[dto.UpdateProject]

type updateProjectHandler struct {
	projectRepository project.Repository
}

func NewUpdateProjectHandler(
	projectRepository project.Repository,
	logger log.Logger,
	metricsClient decorator.MetricsClient,
) UpdateProjectHandler {
	return decorator.ApplyCommandDecorators[dto.UpdateProject](
		updateProjectHandler{
			projectRepository: projectRepository,
		},
		logger,
		metricsClient,
	)
}

func (c updateProjectHandler) Handle(ctx context.Context, cmd dto.UpdateProject) error {
	//TODO implement me
	panic("implement me")
}
