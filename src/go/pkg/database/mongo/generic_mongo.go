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

type Model interface {
	GetID() string
}

type BaseModel struct {
	ID string `bson:"_id"`
}

func NewBaseModel(id string) BaseModel {
	return BaseModel{
		ID: id,
	}
}

func (b BaseModel) GetID() string {
	return b.ID
}

type DataStore[M Model] struct {
	mongoRepository mongodb.MongoRepository
	collectionName  string
}

func NewDataStore[M Model](
	mongoRepository mongodb.MongoRepository,
	collectionName string,
) *DataStore[M] {
	return &DataStore[M]{
		mongoRepository: mongoRepository,
		collectionName:  collectionName,
	}
}

func (m *DataStore[M]) Save(ctx context.Context, model M) error {
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).InsertOne(ctx, model)
	if err != nil {
		return errors.Wrap(err, "failed to save model")
	}
	return nil
}

func (m *DataStore[M]) SaveMany(ctx context.Context, models []M) error {
	docs := make([]interface{}, len(models))
	for i, model := range models {
		docs[i] = model
	}
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).InsertMany(ctx, docs)
	if err != nil {
		return errors.Wrap(err, "failed to save multiple models")
	}
	return nil
}

func (m *DataStore[M]) Get(ctx context.Context, id string) (*M, error) {
	var model M
	filter := bson.M{"_id": id}
	err := m.mongoRepository.Database(ctx).Collection(m.collectionName).FindOne(ctx, filter).Decode(&model)
	if stderrors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.NotFoundError(mongo.ErrNoDocuments)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get model")
	}
	return &model, nil
}

func (m *DataStore[M]) Update(ctx context.Context, model M) error {
	filter := bson.M{"_id": model.GetID()}
	update := bson.M{"$set": model}
	opts := options.Update().SetUpsert(true)
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errors.Wrap(err, "failed to update model")
	}
	return nil
}

func (m *DataStore[M]) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "failed to delete model")
	}
	return nil
}

func (m *DataStore[M]) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]M, error) {
	cursor, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).Find(ctx, filter, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find models")
	}
	defer cursor.Close(ctx)
	var results []M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.Wrap(err, "failed to decode models")
	}
	return results, nil
}

func (m *DataStore[M]) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := m.mongoRepository.Database(ctx).Collection(m.collectionName).CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count models")
	}
	return count, nil
}
