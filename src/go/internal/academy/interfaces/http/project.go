package http

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	openapi "github.com/guodongq/quickstart/pkg/api/codegen"
	"go.opentelemetry.io/otel"
)

func (h *HttpServer) ProjectsList(ctx context.Context, request openapi.ProjectsListRequestObject) (openapi.ProjectsListResponseObject, error) {
	ctx, span := otel.Tracer("Projects").Start(ctx, "HttpServer::ProjectsList")
	defer span.End()

	panic("implement me")
}

func (h *HttpServer) ProjectsCreate(ctx context.Context, request openapi.ProjectsCreateRequestObject) (openapi.ProjectsCreateResponseObject, error) {
	ctx, span := otel.Tracer("Projects").Start(ctx, "HttpServer::ProjectsCreate")
	defer span.End()

	model, err := dto.NewCreateProject(ctx, request)
	if err != nil {
		return nil, err
	}

	if err := h.app.Commands.CreateProject.Handle(ctx, model); err != nil {
		return nil, err
	}

	project, err := h.app.Queries.GetProjectByID.Handle(ctx, dto.GetProjectByID{ID: model.ID})
	if err != nil {
		return nil, err
	}

	return openapi.ProjectsCreate201JSONResponse{
		ExistingProjectJSONResponse: openapi.ExistingProjectJSONResponse(project.ToAPI()),
	}, nil
}

func (h *HttpServer) ProjectsDelete(ctx context.Context, request openapi.ProjectsDeleteRequestObject) (openapi.ProjectsDeleteResponseObject, error) {
	panic("implement me")
}

func (h *HttpServer) ProjectsGet(ctx context.Context, request openapi.ProjectsGetRequestObject) (openapi.ProjectsGetResponseObject, error) {
	panic("implement me")
}

func (h *HttpServer) ProjectsPatch(ctx context.Context, request openapi.ProjectsPatchRequestObject) (openapi.ProjectsPatchResponseObject, error) {
	panic("implement me")
}

func (h *HttpServer) ProjectsUpdate(ctx context.Context, request openapi.ProjectsUpdateRequestObject) (openapi.ProjectsUpdateResponseObject, error) {
	panic("implement me")
}
