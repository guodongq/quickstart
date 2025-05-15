package mongo

import (
	"context"
	"github.com/guodongq/quickstart/internal/academy/domain/environment"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
)

type EnvironmentPO struct {
	BasePO `bson:",inline"`
	Name   string `bson:"name"`
}

func fromEnvironment(entity *environment.Environment) EnvironmentPO {
	return EnvironmentPO{}
}

func (e EnvironmentPO) toEnvironment() *environment.Environment {
	return &environment.Environment{}
}

type environmentRepository struct {
	dataStore *MongoRepository[EnvironmentPO]
}

func NewEnvironmentRepository(repoProvider mongodb.MongoRepository) environment.Repository {
	return &environmentRepository{
		dataStore: NewMongoDataStore[EnvironmentPO](
			repoProvider,
			collectionProject,
		),
	}
}

func (p environmentRepository) CreateEnvironment(ctx context.Context, entity *environment.Environment) error {
	//TODO implement me
	panic("implement me")
}
