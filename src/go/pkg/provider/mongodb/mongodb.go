package mongodb

import (
	"context"
	"time"

	"github.com/guodongq/quickstart/pkg/provider/app"
	"github.com/guodongq/quickstart/pkg/provider/probes"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type MongoDB struct {
	provider.AbstractProvider

	options        MongoDBOptions
	appProvider    *app.App
	probesProvider *probes.Probes
	Client         *mongo.Client
	Database       *mongo.Database
}

func New(appProvider *app.App, probesProvider *probes.Probes, optionFuncs ...func(*MongoDBOptions)) *MongoDB {
	defaultOptions := getDefaultMongoDBOptions()
	options := &defaultOptions

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &MongoDB{
		appProvider:    appProvider,
		probesProvider: probesProvider,
		options:        *options,
	}
}

func (p *MongoDB) Init() error {
	opts := mongoOptions.Client()
	opts.ApplyURI(p.options.URI)
	opts.SetConnectTimeout(p.options.Timeout)
	opts.SetHeartbeatInterval(p.options.HeartbeatInterval)
	opts.SetMaxPoolSize(p.options.MaxPoolSize)
	opts.SetMinPoolSize(p.options.MinPoolSize)
	opts.SetMaxConnIdleTime(p.options.MaxConnIdleTime)
	opts.SetReadPreference(readpref.Primary())
	opts.SetWriteConcern(writeconcern.Majority())
	if p.appProvider != nil {
		opts.SetAppName(p.appProvider.Name())
	}

	uri, err := connstring.Parse(p.options.URI)
	if err != nil {
		return err
	}
	if uri.Database != "" {
		p.options.Database = uri.Database
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
	defer cancel()

	logEntry := logger.WithField("address", p.options.URI).WithField("time_out", p.options.Timeout.String())

	logEntry.Debug("Connecting to MongoDB server...")

	// Create Client and connect to MongoDB.
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB connection failed")
		return err
	}

	// Check connection by pinging.
	err = client.Ping(ctx, nil)
	if err != nil {
		logEntry.WithError(err).Error("MongoDB ping failed")
		return err
	}

	p.Client = client
	p.Database = client.Database(p.options.Database)

	// Add live probes if possible.
	if p.probesProvider != nil {
		p.probesProvider.AddLivenessProbes(p.livenessProbe)
	}

	return nil
}

// Close to connection with the MongoDB server.
func (p *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
	defer cancel()

	err := p.Client.Disconnect(ctx)
	if err != nil {
		logger.WithError(err).Info("MongoDB disconnecting failed")
		return err
	}

	return p.AbstractProvider.Close()
}

func (p *MongoDB) livenessProbe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Client.Ping(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("MongoDB liveness probe failed")
		return err
	}

	return nil
}
