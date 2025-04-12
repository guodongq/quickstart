package rabbitmq

import (
	logger "github.com/guodongq/quickstart/pkg/util/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func getDefaultConsumerOptions() ConsumerOptions {
	return ConsumerOptions{
		RabbitConsumerOptions: RabbitConsumerOptions{
			Name:      "",
			AutoAck:   false,
			Exclusive: false,
			NoWait:    false,
			NoLocal:   false,
			Args:      Table{},
		},
		QueueOptions: QueueOptions{
			Name:       "",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Passive:    false,
			Args:       Table{},
			Declare:    true,
		},
		ExchangeOptions: []ExchangeOptions{},
		Concurrency:     1,
		QOSPrefetch:     1,
		QOSGlobal:       false,
		Handler: func(d Delivery) (action Action) {
			logger.Infof("Delivery: %v", d)
			return Ack
		},
	}
}

func getDefaultExchangeOptions() ExchangeOptions {
	return ExchangeOptions{
		Name:       "",
		Kind:       amqp.ExchangeDirect,
		Durable:    false,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Passive:    false,
		Args:       Table{},
		Declare:    true,
		Bindings:   []Binding{},
	}
}

func getDefaultBindingOptions() BindingOptions {
	return BindingOptions{
		NoWait:  false,
		Args:    Table{},
		Declare: true,
	}
}

type ConsumerOptions struct {
	RabbitConsumerOptions RabbitConsumerOptions
	QueueOptions          QueueOptions
	ExchangeOptions       []ExchangeOptions
	Concurrency           int
	QOSPrefetch           int
	QOSGlobal             bool
	Handler               Handler
}

type RabbitConsumerOptions struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoWait    bool
	NoLocal   bool
	Args      Table
}

type QueueOptions struct {
	Name          string
	TemporaryName string // if the Name is empty, the queue will be named with the result of QueueDeclare
	Durable       bool
	AutoDelete    bool
	Exclusive     bool
	NoWait        bool
	Passive       bool // if false, a missing queue will be created on the server
	Args          Table
	Declare       bool
}

type Binding struct {
	RoutingKey string
	BindingOptions
}

type BindingOptions struct {
	NoWait  bool
	Args    Table
	Declare bool
}

func WithConsumerOptionsQueueName(name string) func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.Name = name
	}
}

func WithConsumerOptionsQueueDurable() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.Durable = true
	}
}

func WithConsumerOptionsQueueAutoDelete() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.AutoDelete = true
	}
}

func WithConsumerOptionsQueueExclusive() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.Exclusive = true
	}
}

func WithConsumerOptionsQueueNoWait() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.NoWait = true
	}
}

func WithConsumerOptionsQueuePassive() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.Passive = true
	}
}

func WithConsumerOptionsQueueArgs(args Table) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.Args = args
	}
}

func WithConsumerOptionsQueueNoDeclare() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QueueOptions.Declare = false
	}
}

func ensureExchangeOptions(options *ConsumerOptions) {
	if len(options.ExchangeOptions) == 0 {
		options.ExchangeOptions = append(options.ExchangeOptions, getDefaultExchangeOptions())
	}
}

// WithConsumerOptionsExchangeName sets the exchange name
func WithConsumerOptionsExchangeName(name string) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Name = name
	}
}

// WithConsumerOptionsExchangeKind ensures the queue is a durable queue
func WithConsumerOptionsExchangeKind(kind string) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Kind = kind
	}
}

// WithConsumerOptionsExchangeDurable ensures the exchange is a durable exchange
func WithConsumerOptionsExchangeDurable() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Durable = true
	}
}

// WithConsumerOptionsExchangeAutoDelete ensures the exchange is an auto-delete exchange
func WithConsumerOptionsExchangeAutoDelete() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].AutoDelete = true
	}
}

// WithConsumerOptionsExchangeInternal ensures the exchange is an internal exchange
func WithConsumerOptionsExchangeInternal() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Internal = true
	}
}

// WithConsumerOptionsExchangeNoWait ensures the exchange is a no-wait exchange
func WithConsumerOptionsExchangeNoWait() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].NoWait = true
	}
}

// WithConsumerOptionsExchangeDeclare stops this library from declaring the exchanges existance
func WithConsumerOptionsExchangeDeclare() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Declare = true
	}
}

// WithConsumerOptionsExchangePassive ensures the exchange is a passive exchange
func WithConsumerOptionsExchangePassive() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Passive = true
	}
}

// WithConsumerOptionsExchangeArgs adds optional args to the exchange
func WithConsumerOptionsExchangeArgs(args Table) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Args = args
	}
}

// WithConsumerOptionsRoutingKey binds the queue to a routing key with the default binding options
func WithConsumerOptionsRoutingKey(routingKey string) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Bindings = append(options.ExchangeOptions[0].Bindings, Binding{
			RoutingKey:     routingKey,
			BindingOptions: getDefaultBindingOptions(),
		})
	}
}

// WithConsumerOptionsBinding adds a new binding to the queue which allows you to set the binding options
// on a per-binding basis. Keep in mind that everything in the BindingOptions struct will default to
// the zero value. If you want to declare your bindings for example, be sure to set Declare=true
func WithConsumerOptionsBinding(binding Binding) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		ensureExchangeOptions(options)
		options.ExchangeOptions[0].Bindings = append(options.ExchangeOptions[0].Bindings, binding)
	}
}

// WithConsumerOptionsExchangeOptions adds a new exchange to the consumer, this should probably only be
// used if you want to to consume from multiple exchanges on the same consumer
func WithConsumerOptionsExchangeOptions(exchangeOptions ExchangeOptions) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.ExchangeOptions = append(options.ExchangeOptions, exchangeOptions)
	}
}

// WithConsumerOptionsConcurrency returns a function that sets the concurrency, which means that
// many goroutines will be spawned to run the provided handler on messages
func WithConsumerOptionsConcurrency(concurrency int) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.Concurrency = concurrency
	}
}

// WithConsumerOptionsConsumerName returns a function that sets the name on the server of this consumer
// if unset a random name will be given
func WithConsumerOptionsConsumerName(consumerName string) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.RabbitConsumerOptions.Name = consumerName
	}
}

// WithConsumerOptionsConsumerAutoAck returns a function that sets the auto acknowledge property on the server of this consumer
// if unset the default will be used (false)
func WithConsumerOptionsConsumerAutoAck(autoAck bool) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.RabbitConsumerOptions.AutoAck = autoAck
	}
}

// WithConsumerOptionsConsumerExclusive sets the consumer to exclusive, which means
// the server will ensure that this is the sole consumer
// from this queue. When exclusive is false, the server will fairly distribute
// deliveries across multiple consumers.
func WithConsumerOptionsConsumerExclusive() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.RabbitConsumerOptions.Exclusive = true
	}
}

// WithConsumerOptionsConsumerNoWait sets the consumer to nowait, which means
// it does not wait for the server to confirm the request and
// immediately begin deliveries. If it is not possible to consume, a channel
// exception will be raised and the channel will be closed.
func WithConsumerOptionsConsumerNoWait() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.RabbitConsumerOptions.NoWait = true
	}
}

func WithConsumerOptionsConsumerNoLocal() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.RabbitConsumerOptions.NoLocal = true
	}
}

func WithConsumerOptionsConsumerArgs(args Table) func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.RabbitConsumerOptions.Args = args
	}
}

// WithConsumerOptionsQOSPrefetch returns a function that sets the prefetch count, which means that
// many messages will be fetched from the server in advance to help with throughput.
// This doesn't affect the handler, messages are still processed one at a time.
func WithConsumerOptionsQOSPrefetch(prefetchCount int) func(*ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QOSPrefetch = prefetchCount
	}
}

// WithConsumerOptionsQOSGlobal sets the qos on the channel to global, which means
// these QOS settings apply to ALL existing and future
// consumers on all channels on the same connection
func WithConsumerOptionsQOSGlobal() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.QOSGlobal = true
	}
}

// WithConsumerOptionsQueueQuorum sets the queue a quorum type, which means
// multiple nodes in the cluster will have the messages distributed amongst them
// for higher reliability
func WithConsumerOptionsQueueQuorum() func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		if options.QueueOptions.Args == nil {
			options.QueueOptions.Args = Table{}
		}

		options.QueueOptions.Args["x-queue-type"] = "quorum"
	}
}

// WithConsumerOptionsLogger sets handler to a custom interface.
// Use WithConsumerOptionsHandler to register consumer handler
func WithConsumerOptionsHandler(handler Handler) func(options *ConsumerOptions) {
	return func(options *ConsumerOptions) {
		options.Handler = handler
	}
}
