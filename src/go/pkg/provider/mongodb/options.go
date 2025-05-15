package mongodb

import (
	"time"
)

type MongoDBOptions struct {
	URI               string
	Timeout           time.Duration
	MaxPoolSize       uint64
	MinPoolSize       uint64
	MaxConnIdleTime   time.Duration
	HeartbeatInterval time.Duration
	Database          string
}

func getDefaultMongoDBOptions() MongoDBOptions {
	return MongoDBOptions{
		URI:               "mongodb://root:root@localhost:27017/?directConnection=true",
		Timeout:           time.Second * 20,
		MaxPoolSize:       16,
		MinPoolSize:       1,
		MaxConnIdleTime:   time.Second * 30,
		HeartbeatInterval: time.Second * 15,
	}
}

func (o *MongoDBOptions) MergeIn(opts ...func(*MongoDBOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithMongoDBOptionsURI(uri string) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.URI = uri
	}
}

func WithMongoDBOptionsTimeout(timeout time.Duration) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.Timeout = timeout
	}
}

func WithMongoDBOptionsMaxPoolSize(maxPoolSize uint64) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.MaxPoolSize = maxPoolSize
	}
}

func WithMongoDBOptionsMinPoolSize(minPoolSize uint64) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.MinPoolSize = minPoolSize
	}
}

func WithMongoDBOptionsMaxConnIdleTime(maxConnIdleTime time.Duration) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.MaxConnIdleTime = maxConnIdleTime
	}
}

func WithMongoDBOptionsHeartbeatInterval(heartbeatInterval time.Duration) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.HeartbeatInterval = heartbeatInterval
	}
}

func WithMongoDBOptionsDatabase(database string) func(*MongoDBOptions) {
	return func(o *MongoDBOptions) {
		o.Database = database
	}
}
