package rabbitmq

import (
	"context"
	"fmt"

	"github.com/guodongq/quickstart/pkg/util/provider"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

type Action int

type Handler func(d Delivery) (action Action)

const (
	// Ack default ack this msg after you have successfully processed this delivery.
	Ack Action = iota
	// NackDiscard the message will be dropped or delivered to a server configured dead-letter queue.
	NackDiscard
	// NackRequeue deliver this message to a different consumer.
	NackRequeue
	// Message acknowledgement is left to the user using the msg.Ack() method
	Manual
)

type Consumer struct {
	provider.AbstractRunProvider
	connection *Connection
	options    ConsumerOptions
}

type Delivery struct {
	amqp.Delivery
}

func NewConsumer(
	connection *Connection,
	optionFuncs ...func(*ConsumerOptions),
) *Consumer {
	defaultOptions := getDefaultConsumerOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Consumer{
		connection: connection,
		options:    *options,
	}
}

func (p *Consumer) QueueName() string {
	if p.options.QueueOptions.Name != "" {
		return p.options.QueueOptions.Name
	}
	return p.options.QueueOptions.TemporaryName
}

func (p *Consumer) Init() (err error) {
	for _, exchangeOption := range p.options.ExchangeOptions {
		err = declareExchange(p.connection, exchangeOption)
		if err != nil {
			return fmt.Errorf("declare exchange failed: %w", err)
		}
	}

	err = declareQueue(p.connection, &p.options.QueueOptions)
	if err != nil {
		return fmt.Errorf("declare queue failed: %w", err)
	}

	err = declareBindings(p.connection, p.options)
	if err != nil {
		return fmt.Errorf("declare bindings failed: %w", err)
	}
	return nil
}

func (p *Consumer) Run() error {
	err := p.connection.ch.Qos(
		p.options.QOSPrefetch,
		0,
		p.options.QOSGlobal,
	)
	if err != nil {
		return fmt.Errorf("declare qos failed: %w", err)
	}

	msgs, err := p.connection.ch.Consume(
		p.options.QueueOptions.Name,
		p.options.RabbitConsumerOptions.Name,
		p.options.RabbitConsumerOptions.AutoAck,
		p.options.RabbitConsumerOptions.Exclusive,
		false, // no-local is not supported by RabbitMQ
		p.options.RabbitConsumerOptions.NoWait,
		tableToAMQPTable(p.options.RabbitConsumerOptions.Args),
	)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(context.Background())
	for i := 0; i < p.options.Concurrency; i++ {
		g.Go(func() error {
			return handlerGoroutine(gCtx, msgs, p.options, p.options.Handler)
		})
	}
	return g.Wait()
}

func handlerGoroutine(ctx context.Context, msgs <-chan amqp.Delivery, consumerOptions ConsumerOptions, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("RabbitMQ channel closed while consuming messages")
			}

			if consumerOptions.RabbitConsumerOptions.AutoAck {
				handler(Delivery{msg})
				continue
			}

			switch handler(Delivery{msg}) {
			case Ack:
				err := msg.Ack(false)
				if err != nil {
					return fmt.Errorf("can't ack message: %v", err)
				}
			case NackDiscard:
				err := msg.Nack(false, false)
				if err != nil {
					return fmt.Errorf("can't nack message: %v", err)
				}
			case NackRequeue:
				err := msg.Nack(false, true)
				if err != nil {
					return fmt.Errorf("can't nack message: %v", err)
				}
			case Manual:
				// do nothing here
			}
		}
	}
}
