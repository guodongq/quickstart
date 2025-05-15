package mongo

import (
	"context"
	stderrors "errors"
	"github.com/guodongq/quickstart/pkg/errors"
	"github.com/guodongq/quickstart/pkg/provider/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionProject     = "project"
	collectionEnvironment = "environment"
)

type Identifiable interface {
	GetID() string
}

type BasePO struct {
	ID string `bson:"_id"`
}

func (b BasePO) GetID() string {
	return b.ID
}

type MongoRepository[PO Identifiable] struct {
	repoProvider   mongodb.MongoRepository
	collectionName string
}

func NewMongoDataStore[PO Identifiable](
	repoProvider mongodb.MongoRepository,
	collectionName string,
) *MongoRepository[PO] {
	return &MongoRepository[PO]{
		repoProvider:   repoProvider,
		collectionName: collectionName,
	}
}

func (m *MongoRepository[PO]) Save(ctx context.Context, model PO) error {
	_, err := m.repoProvider.Database(ctx).Collection(m.collectionName).InsertOne(ctx, model)
	return err
}

func (m *MongoRepository[PO]) SaveMany(ctx context.Context, models []PO) error {
	docs := make([]any, len(models))
	for i, model := range models {
		docs[i] = model
	}
	_, err := m.repoProvider.Database(ctx).Collection(m.collectionName).InsertMany(ctx, docs)
	return err
}

func (m *MongoRepository[PO]) Get(ctx context.Context, id string) (*PO, error) {
	var model PO
	filter := bson.M{"_id": id}
	err := m.repoProvider.Database(ctx).Collection(m.collectionName).FindOne(ctx, filter).Decode(&model)
	if stderrors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.NotFoundError(mongo.ErrNoDocuments)
	}

	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (m *MongoRepository[PO]) Update(ctx context.Context, model PO) error {
	filter := bson.M{"_id": model.GetID()}
	update := bson.M{"$set": model}
	opts := options.Update().SetUpsert(true)
	_, err := m.repoProvider.Database(ctx).Collection(m.collectionName).UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *MongoRepository[PO]) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := m.repoProvider.Database(ctx).Collection(m.collectionName).DeleteOne(ctx, filter)
	return err
}

func (m *MongoRepository[PO]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]PO, error) {
	cursor, err := m.repoProvider.Database(ctx).Collection(m.collectionName).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []PO
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoRepository[PO]) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := m.repoProvider.Database(ctx).Collection(m.collectionName).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
