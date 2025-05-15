package app

import (
	"github.com/guodongq/quickstart/internal/academy/app/commands"
	"github.com/guodongq/quickstart/internal/academy/app/event"
	"github.com/guodongq/quickstart/internal/academy/app/queries"
	"github.com/guodongq/quickstart/internal/academy/domain/environment"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/bus"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateProject commands.CreateProjectHandler
}

type Queries struct {
	GetProjectByID queries.GetProjectByIDHandler
}

// NewApp creates a new application instance
func NewApp(
	bus bus.EventBus,
) *App {
	var (
		createProject   commands.CreateProjectHandler
		getProjectByID  queries.GetProjectByIDHandler
		projectRepo     project.Repository
		environmentRepo environment.Repository
	)

	// register event handlers
	for _, eventRegistry := range []event.EventRegistry{
		event.NewProjectEventRegistry(bus, projectRepo),
		event.NewEnvironmentEventRegistry(bus, environmentRepo),
	} {
		eventRegistry.RegisterHandlers()
	}

	return &App{
		Commands: Commands{
			CreateProject: createProject,
		},
		Queries: Queries{
			GetProjectByID: getProjectByID,
		},
	}
}
