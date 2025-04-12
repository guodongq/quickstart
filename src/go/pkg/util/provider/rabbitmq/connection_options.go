package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionOptions struct {
	URL    string
	Config Config
}

func getDefaultConnectionOptions() ConnectionOptions {
	return ConnectionOptions{
		Config: Config{
			Heartbeat:  time.Second * 10,
			Locale:     "en_US",
			Properties: amqp.Table{},
		},
	}
}

func WithConnectionOptionsURL(url string) func(options *ConnectionOptions) {
	return func(options *ConnectionOptions) {
		options.URL = url
	}
}

func WithConnectionOptionsConfig(config Config) func(options *ConnectionOptions) {
	return func(options *ConnectionOptions) {
		options.Config = config
	}
}

func WithConnectionOptionsConfigHearBeat(heartbeat time.Duration) func(options *ConnectionOptions) {
	return func(options *ConnectionOptions) {
		options.Config.Heartbeat = heartbeat
	}
}

func WithConnectionOptionsConfigProperties(properties Table) func(options *ConnectionOptions) {
	return func(options *ConnectionOptions) {
		options.Config.Properties = amqp.Table(properties)
	}
}

func WithConnectionOptionsConfigVhost(vhost string) func(options *ConnectionOptions) {
	return func(options *ConnectionOptions) {
		options.Config.Vhost = vhost
	}
}

func WithConnectionOptionsConfigPropertiesConnectionName(connectionName string) func(options *ConnectionOptions) {
	return func(options *ConnectionOptions) {
		options.Config.Properties["connection_name"] = connectionName
	}
}
