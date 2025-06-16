package queries

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/decorator"
	"github.com/guodongq/quickstart/pkg/log"
)

type GetProjectByIDHandler decorator.QueryHandler[dto.GetProjectByID, dto.Project]

type getProjectByIDHandler struct {
	projectRepository project.Repository
}

func NewGetProjectByIDHandler(
	projectRepository project.Repository,
	logger log.Logger,
	metricsClient decorator.MetricsClient,
) GetProjectByIDHandler {
	return decorator.ApplyQueryDecorators[dto.GetProjectByID, dto.Project](
		getProjectByIDHandler{
			projectRepository: projectRepository,
		},
		logger,
		metricsClient,
	)
}

func (c getProjectByIDHandler) Handle(ctx context.Context, query dto.GetProjectByID) (dto.Project, error) {
	resp, err := c.projectRepository.GetProjectByID(ctx, query.ID)
	if err != nil {
		return dto.Project{}, err
	}

	return dto.NewProjectDto(*resp), nil
}
