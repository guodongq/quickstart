package project

import (
	"context"
	"github.com/guodongq/quickstart/pkg/idgen"
)

type Repository interface {
	CreateProject(ctx context.Context, entity *Project) error
	GetProjectByID(ctx context.Context, id idgen.Generator) (*Project, error)
	UpdateProjectByID(
		ctx context.Context,
		id idgen.Generator,
		updateFn func(entity *Project) (*Project, error),
	) error
	DeleteProjectByID(ctx context.Context, id idgen.Generator) error
}
