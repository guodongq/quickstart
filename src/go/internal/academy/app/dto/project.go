package dto

import (
	"context"
	"github.com/google/uuid"
	openapi "github.com/guodongq/quickstart/pkg/api/codegen"
	"github.com/guodongq/quickstart/pkg/errors"
	"github.com/guodongq/quickstart/pkg/utils"
)

type Project struct {
}

func (p Project) ToAPI() openapi.Project {
	return openapi.Project{}
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
