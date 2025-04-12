package mongodb

import (
	"context"
	"github.com/google/uuid"
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoRepository interface {
	provider.Provider
	Client() *mongo.Client
	Database(ctx context.Context) *mongo.Database
	RunTransaction(ctx context.Context, txnFn func(sctx mongo.SessionContext) error) error
}

type DatabaseNameGetter func() string

func WithDatabaseNameGetter(databaseNameGetter DatabaseNameGetter) func(*mongoRepository) {
	return func(m *mongoRepository) {
		m.databaseNameGetter = databaseNameGetter
	}
}

type databaseNameCtx struct{}

func GetDatabaseName(ctx context.Context) string {
	databaseName, exist := ctx.Value(databaseNameCtx{}).(string)
	if exist {
		return databaseName
	}
	return ""
}

func WithDatabaseName(ctx context.Context, databaseName string) context.Context {
	return context.WithValue(ctx, databaseNameCtx{}, databaseName)
}

type mongoRepository struct {
	provider.AbstractProvider

	mongoProvider      *MongoDB
	databaseNameGetter DatabaseNameGetter
}

func NewMongoRepository(mongoProvider *MongoDB, optionFuncs ...func(*mongoRepository)) MongoRepository {
	mongoRepository := &mongoRepository{
		mongoProvider: mongoProvider,
	}

	for _, optionFunc := range optionFuncs {
		optionFunc(mongoRepository)
	}

	return mongoRepository
}

func (m mongoRepository) Client() *mongo.Client {
	return m.mongoProvider.Client
}

func (m mongoRepository) Database(ctx context.Context) *mongo.Database {
	if m.databaseNameGetter != nil {
		return m.mongoProvider.Client.Database(m.databaseNameGetter())
	}

	if databaseName := GetDatabaseName(ctx); databaseName != "" {
		return m.mongoProvider.Client.Database(databaseName)
	}

	return m.mongoProvider.Database
}

func (m mongoRepository) RunTransaction(ctx context.Context, txnFn func(sctx mongo.SessionContext) error) error {
	return m.mongoProvider.Client.UseSessionWithOptions(ctx,
		options.Session().SetDefaultReadPreference(readpref.Primary()), func(sessionContext mongo.SessionContext) error {
			return m.runTransactionWithRetry(sessionContext, txnFn)
		})
}

func (m mongoRepository) runTransactionWithRetry(sctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
	txID := uuid.NewString()
	logger.Info("Begin transaction: ", txID)
	for {
		// start transaction
		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			logger.Info("Start transaction: %s. Error: %s\n", txID, err)
			return err
		}

		err = txnFn(sctx)
		if err != nil {
			logger.Info("Abort transaction: %s. Error: %s\n", txID, err)
			if abortErr := sctx.AbortTransaction(sctx); abortErr != nil {
				logger.Info("Failed to abort transaction: ", txID)
			}
			return err
		}

		err = m.commitWithRetry(sctx)
		switch e := err.(type) {
		case nil:
			logger.Info("End transaction: %s, successful. \n", txID)
			return nil
		case mongo.CommandError:
			// If transient error, retry the whole transaction
			if e.HasErrorLabel("TransientTransactionError") {
				logger.Info("TransientTransactionError: %s, retrying transaction...\n", txID)
				continue
			}
			return e
		default:
			logger.Info("End transaction: %s. Error: %s\n", txID, err)
			return e
		}
	}
}

func (m mongoRepository) commitWithRetry(sctx mongo.SessionContext) error {
	for {
		err := sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			return nil
		case mongo.CommandError:
			// Can retry commit
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				continue
			}
			return e
		default:
			return e
		}
	}
}
