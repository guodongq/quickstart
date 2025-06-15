package http

import (
	"github.com/guodongq/quickstart/internal/academy/app"
	openapi "github.com/guodongq/quickstart/pkg/api/codegen"
)

type HttpServer struct {
	openapi.StrictServerInterface
	app *app.App
}

func NewHttpServer(app *app.App) *HttpServer {
	return &HttpServer{app: app}
}
