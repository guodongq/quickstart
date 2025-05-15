package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherOptions struct {
	ExchangeOptions ExchangeOptions
	ConfirmMode     bool
}

func getDefaultPublisherOptions() PublisherOptions {
	return PublisherOptions{
		ExchangeOptions: ExchangeOptions{
			Name:       "",
			Kind:       amqp.ExchangeDirect,
			Durable:    false,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Passive:    false,
			Args:       Table{},
			Declare:    false,
		},
		ConfirmMode: false,
	}
}

// WithPublisherOptionsExchangeName sets the exchange name
func WithPublisherOptionsExchangeName(name string) func(*PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Name = name
	}
}

// WithPublisherOptionsExchangeKind ensures the queue is a durable queue
func WithPublisherOptionsExchangeKind(kind string) func(*PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Kind = kind
	}
}

// WithPublisherOptionsExchangeDurable ensures the exchange is a durable exchange
func WithPublisherOptionsExchangeDurable() func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Durable = true
	}
}

// WithPublisherOptionsExchangeAutoDelete ensures the exchange is an auto-delete exchange
func WithPublisherOptionsExchangeAutoDelete() func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.AutoDelete = true
	}
}

// WithPublisherOptionsExchangeInternal ensures the exchange is an internal exchange
func WithPublisherOptionsExchangeInternal() func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Internal = true
	}
}

// WithPublisherOptionsExchangeNoWait ensures the exchange is a no-wait exchange
func WithPublisherOptionsExchangeNoWait() func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.NoWait = true
	}
}

// WithPublisherOptionsExchangeDeclare will create the exchange if it doesn't exist
func WithPublisherOptionsExchangeDeclare() func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Declare = true
	}
}

// WithPublisherOptionsExchangePassive ensures the exchange is a passive exchange
func WithPublisherOptionsExchangePassive() func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Passive = true
	}
}

// WithPublisherOptionsExchangeArgs adds optional args to the exchange
func WithPublisherOptionsExchangeArgs(args Table) func(*PublisherOptions) {
	return func(options *PublisherOptions) {
		options.ExchangeOptions.Args = args
	}
}

// WithPublisherOptionsConfirm enables confirm mode on the connection
// this is required if publisher confirmations should be used
func WithPublisherOptionsConfirm(options *PublisherOptions) {
	options.ConfirmMode = true
}
