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

type Identifiable interface {
	GetID() string
}

type BaseModel struct {
	ID string `bson:"_id"`
}

func (b BaseModel) GetID() string {
	return b.ID
}

type DataStore[P Identifiable] struct {
	mongoRepository mongodb.MongoRepository
	collectionName  string
}

func NewDataStore[P Identifiable](
	mongoRepository mongodb.MongoRepository,
	collectionName string,
) *DataStore[P] {
	return &DataStore[P]{
		mongoRepository: mongoRepository,
		collectionName:  collectionName,
	}
}

func (m *DataStore[P]) Save(ctx context.Context, model P) error {
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).InsertOne(ctx, model)
	return err
}

func (m *DataStore[P]) SaveMany(ctx context.Context, models []P) error {
	docs := make([]any, len(models))
	for i, model := range models {
		docs[i] = model
	}
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).InsertMany(ctx, docs)
	return err
}

func (m *DataStore[P]) Get(ctx context.Context, id string) (*P, error) {
	var model P
	filter := bson.M{"_id": id}
	err := m.mongoRepository.Database(ctx).Collection(m.collectionName).FindOne(ctx, filter).Decode(&model)
	if stderrors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.NotFoundError(mongo.ErrNoDocuments)
	}

	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (m *DataStore[P]) Update(ctx context.Context, model P) error {
	filter := bson.M{"_id": model.GetID()}
	update := bson.M{"$set": model}
	opts := options.Update().SetUpsert(true)
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *DataStore[P]) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).DeleteOne(ctx, filter)
	return err
}

func (m *DataStore[P]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]P, error) {
	cursor, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []P
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *DataStore[P]) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
