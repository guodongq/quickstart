package app

import (
	"github.com/guodongq/quickstart/internal/academy/app/commands"
	"github.com/guodongq/quickstart/internal/academy/app/events"
	"github.com/guodongq/quickstart/internal/academy/app/queries"
	"github.com/guodongq/quickstart/internal/academy/infra/persistence/mongo"
	"github.com/guodongq/quickstart/pkg/bus"
	"github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/metrics"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateProject commands.CreateProjectHandler
	UpdateProject commands.UpdateProjectHandler
	DeleteProject commands.DeleteProjectHandler
}

type Queries struct {
	GetProjectByID queries.GetProjectByIDHandler
}

// NewApplication creates a new application instance
func NewApplication(
	mongoRepository mongodb.MongoRepository,
) *App {
	var (
		eventBus      = bus.New()
		metricsClient = metrics.New()
		logger        = log.DefaultLogger()
	)

	// project
	var (
		projectRepository = mongo.NewProjectRepository(mongoRepository)
	)

	// environment
	var (
		environmentRepository = mongo.NewEnvironmentRepository(mongoRepository)
	)

	// register event handlers
	events.NewDefaultEventRegistryManager(
		events.NewProjectEventRegistry(eventBus, projectRepository),
		events.NewEnvironmentEventRegistry(eventBus, environmentRepository),
	).Register()

	return &App{
		Commands: Commands{
			CreateProject: commands.NewCreateProjectHandler(projectRepository, logger, metricsClient, eventBus),
			UpdateProject: commands.NewUpdateProjectHandler(projectRepository, logger, metricsClient),
			DeleteProject: commands.NewDeleteProjectHandler(projectRepository, logger, metricsClient),
		},
		Queries: Queries{
			GetProjectByID: queries.NewGetProjectByIDHandler(projectRepository, logger, metricsClient),
		},
	}
}
