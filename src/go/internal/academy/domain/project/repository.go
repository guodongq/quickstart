package project

import (
	"context"
)

type Repository interface {
	CreateProject(ctx context.Context, entity *Project) error
	GetProjectByID(ctx context.Context, id string) (*Project, error)
	UpdateProjectByID(
		ctx context.Context,
		id string,
		updateFn func(entity *Project) (*Project, error),
	) error
	DeleteProjectByID(ctx context.Context, id string) error
}
