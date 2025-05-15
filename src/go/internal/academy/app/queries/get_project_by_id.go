package queries

import (
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	"github.com/guodongq/quickstart/pkg/decorator"
)

type GetProjectByIDHandler decorator.QueryHandler[dto.GetProjectByID, dto.Project]
