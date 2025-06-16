package mongo

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/domain/project"
	"github.com/guodongq/quickstart/pkg/database/mongo"
	"github.com/guodongq/quickstart/pkg/ddd"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
	"github.com/guodongq/quickstart/pkg/types"
)

type ProjectModel struct {
	mongo.BaseModel `bson:",inline"`
	ProjectName     string      `bson:"name"`
	Description     string      `bson:"description"`
	Meta            *types.Meta `bson:"meta"`
}

func fromProject(entity *project.Project) ProjectModel {
	return ProjectModel{
		BaseModel:   mongo.NewBaseModel(entity.Id),
		ProjectName: entity.Name,
		Description: entity.Description,
		Meta:        entity.Meta,
	}
}

func (p ProjectModel) toProject() *project.Project {
	return &project.Project{
		BaseEntity:  ddd.NewBaseEntity(p.ID),
		Name:        p.ProjectName,
		Description: p.Description,
		Meta:        p.Meta,
	}
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

func (p projectRepository) GetProjectByID(ctx context.Context, id string) (*project.Project, error) {
	po, err := p.dataStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return po.toProject(), nil
}

func (p projectRepository) UpdateProjectByID(ctx context.Context, id string, updateFn func(entity *project.Project) (*project.Project, error)) error {
	po, err := p.dataStore.Get(ctx, id)
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

func (p projectRepository) DeleteProjectByID(ctx context.Context, id string) error {
	return p.dataStore.Delete(ctx, id)
}
