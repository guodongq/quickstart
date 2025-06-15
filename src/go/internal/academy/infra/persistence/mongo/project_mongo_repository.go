package mongo

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/database/mongo"
	"github.com/guodongq/quickstart/pkg/idgen"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
)

type ProjectModel struct {
	mongo.BaseModel `bson:",inline"`
	ProjectName     string `bson:"project_name"`
}

func fromProject(entity *project.Project) ProjectModel {
	return ProjectModel{}
}

func (p ProjectModel) toProject() *project.Project {
	return &project.Project{}
}

const collectionProject = "project"

type projectRepository struct {
	dataStore *mongo.DataStore[ProjectModel]
}

func NewProjectRepository(repoProvider mongodb.MongoRepository) project.Repository {
	return &projectRepository{
		dataStore: mongo.NewDataStore[ProjectModel](
			repoProvider,
			collectionProject,
		),
	}
}

func (p projectRepository) CreateProject(ctx context.Context, entity *project.Project) error {
	return p.dataStore.Save(ctx, fromProject(entity))
}

func (p projectRepository) GetProjectByID(ctx context.Context, idGenerator idgen.Generator) (*project.Project, error) {
	po, err := p.dataStore.Get(ctx, idGenerator.Generate())
	if err != nil {
		return nil, err
	}
	return po.toProject(), nil
}

func (p projectRepository) UpdateProjectByID(ctx context.Context, id idgen.Generator, updateFn func(entity *project.Project) (*project.Project, error)) error {
	po, err := p.dataStore.Get(ctx, id.Generate())
	if err != nil {
		return err
	}
	entity := po.toProject()
	updatedEntity, err := updateFn(entity)
	if err != nil {
		return err
	}

	return p.dataStore.Update(ctx, fromProject(updatedEntity))
}

func (p projectRepository) DeleteProjectByID(ctx context.Context, idGenerator idgen.Generator) error {
	return p.dataStore.Delete(ctx, idGenerator.Generate())
}
