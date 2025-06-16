package dto

import (
	"context"
	"github.com/google/uuid"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	openapi "github.com/guodongq/quickstart/pkg/api/codegen"
	"github.com/guodongq/quickstart/pkg/errors"
	"github.com/guodongq/quickstart/pkg/utils"
	"time"
)

type Project struct {
	ID          string
	Name        string
	Description string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
	DeletedBy *string
	Version   int
}

func NewProjectDto(entity project.Project) Project {
	return Project{
		ID:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
		CreatedAt:   *entity.Meta.CreatedAt,
		CreatedBy:   *entity.Meta.CreatedBy,
		UpdatedAt:   *entity.Meta.UpdatedAt,
		UpdatedBy:   *entity.Meta.UpdatedBy,
		DeletedAt:   entity.Meta.DeletedAt,
		DeletedBy:   entity.Meta.DeletedBy,
	}
}

func (p Project) ToAPI() openapi.Project {
	id := uuid.MustParse(p.ID)
	createdBy := uuid.MustParse(p.CreatedBy)
	updatedBy := uuid.MustParse(p.UpdatedBy)
	var deletedBy uuid.UUID
	if p.DeletedBy != nil {
		deletedBy = uuid.MustParse(*p.DeletedBy)
	}
	return openapi.Project{
		Description: &p.Description,
		Id:          &id,
		Meta: &openapi.MetaProperties{
			CreatedAt: &p.CreatedAt,
			CreatedBy: &createdBy,
			DeletedAt: p.DeletedAt,
			DeletedBy: &deletedBy,
			UpdatedAt: &p.UpdatedAt,
			UpdatedBy: &updatedBy,
			Version:   &p.Version,
		},
		Name: &p.Name,
	}
}

type CreateProject struct {
	ID          string
	Name        string `validate:"required"`
	Description string
}

func NewCreateProject(
	ctx context.Context,
	req openapi.ProjectsCreateRequestObject,
) (CreateProject, error) {
	var model = CreateProject{
		ID: uuid.NewString(),
	}

	if err := utils.GetValidator().StructCtx(ctx, req); err != nil {
		return model, errors.BadRequestError(err)
	}

	model.Name = *req.Body.Name
	if req.Body.Description != nil {
		model.Description = *req.Body.Description
	}

	return model, nil
}

type UpdateProject struct {
}

type DeleteProject struct {
}

type GetProjectByID struct {
	ID string
}
