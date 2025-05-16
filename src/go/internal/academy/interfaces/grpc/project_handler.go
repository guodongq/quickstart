package grpc

import "github.com/guodongq/quickstart/internal/academy/app"

type ProjectAPIHandler struct {
	app app.App
}

func NewProjectAPIHandler(app app.App) *ProjectAPIHandler {
	return &ProjectAPIHandler{app: app}
}
