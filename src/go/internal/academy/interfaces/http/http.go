package http

import (
	"github.com/gin-gonic/gin"
	"github.com/guodongq/quickstart/internal/academy/app"
	"github.com/guodongq/quickstart/internal/academy/app/dto"
	genapi "github.com/guodongq/quickstart/pkg/api/genapi/server"
)

type HttpServer struct {
	genapi.ServerInterface
	app *app.App
}

func NewHttpServer(app *app.App) *HttpServer {
	return &HttpServer{app: app}
}

func (h *HttpServer) ProjectsList(c *gin.Context, params genapi.ProjectsListParams) {

}

func (h *HttpServer) ProjectsCreate(c *gin.Context) {
	err := h.app.Commands.CreateProject.Handle(c.Request.Context(), dto.CreateProject{})
	if err != nil {
		c.JSON(500, err)
		return
	}

	//h.app.Queries.GetProjectByID.Handle(c.Request.Context(), projectId.ProjectId)
}

func (h *HttpServer) ProjectsDelete(c *gin.Context, projectId genapi.ProjectId) {

}

func (h *HttpServer) ProjectsGet(c *gin.Context, projectId genapi.ProjectId) {

}

func (h *HttpServer) ProjectsPatch(c *gin.Context, projectId genapi.ProjectId) {

}

func (h *HttpServer) ProjectsUpdate(c *gin.Context, projectId genapi.ProjectId) {

}
