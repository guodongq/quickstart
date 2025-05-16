package app

import (
	"github.com/guodongq/quickstart/internal/academy/app/commands"
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
}

type Queries struct {
	GetProjectByID queries.GetProjectByIDHandler
}

// NewApplication creates a new application instance
func NewApplication(
	mongoRepository mongodb.MongoRepository,
) *App {
	var (
		projectRepository     = mongo.NewProjectRepository(mongoRepository)
		environmentRepository = mongo.NewEnvironmentRepository(mongoRepository)

		eventBus      = bus.New()
		metricsClient = metrics.New()
		logger        = log.DefaultLogger()
	)

	// register event handlers
	for _, eventRegistry := range []EventRegistry{
		NewProjectEventRegistry(eventBus, projectRepository),
		NewEnvironmentEventRegistry(eventBus, environmentRepository),
	} {
		eventRegistry.RegisterHandlers()
	}

	return &App{
		Commands: Commands{
			CreateProject: commands.NewCreateProjectHandler(projectRepository, logger, metricsClient),
		},
		Queries: Queries{
			GetProjectByID: queries.NewGetProjectByIDHandler(projectRepository, logger, metricsClient),
		},
	}
}
