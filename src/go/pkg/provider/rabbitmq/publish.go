package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Transient  uint8 = amqp.Transient
	Persistent uint8 = amqp.Persistent
)

type Return struct {
	amqp.Return
}

type Confirmation struct {
	amqp.Confirmation
}

type Publisher struct {
	provider.AbstractProvider
	connection *Connection

	disablePublishDueToFlow    bool
	disablePublishDueToFlowMux *sync.RWMutex

	disablePublishDueToBlocked    bool
	disablePublishDueToBlockedMux *sync.RWMutex

	handlerMux           *sync.Mutex
	notifyReturnHandler  func(r Return)
	notifyPublishHandler func(p Confirmation)

	options PublisherOptions
}

type PublisherConfirmation []*amqp.DeferredConfirmation

func NewPublisher(
	connection *Connection,
	optionFuncs ...func(*PublisherOptions),
) *Publisher {
	defaultOptions := getDefaultPublisherOptions()
	options := &defaultOptions

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Publisher{
		connection:                    connection,
		options:                       *options,
		disablePublishDueToFlow:       false,
		disablePublishDueToFlowMux:    &sync.RWMutex{},
		disablePublishDueToBlocked:    false,
		disablePublishDueToBlockedMux: &sync.RWMutex{},
		handlerMux:                    &sync.Mutex{},
		notifyReturnHandler:           nil,
		notifyPublishHandler:          nil,
	}
}

func (p *Publisher) Init() error {
	err := p.startup()
	if err != nil {
		return err
	}

	if p.options.ConfirmMode {
		p.NotifyPublish(func(_ Confirmation) {
			// set a blank handler to set the channel in confirm mode
		})
	}

	return nil
}

func (p *Publisher) NotifyPublish(handler func(p Confirmation)) {
	p.handlerMux.Lock()
	shouldStart := p.notifyPublishHandler == nil
	p.notifyPublishHandler = handler
	p.handlerMux.Unlock()

	if shouldStart {
		p.startPublishHandler()
	}
}

func (p *Publisher) startPublishHandler() {
	p.handlerMux.Lock()
	if p.notifyPublishHandler == nil {
		p.handlerMux.Unlock()
		return
	}
	p.handlerMux.Unlock()
	p.connection.ch.Confirm(false)

	go func() {
		confirmationCh := p.connection.ch.NotifyPublish(make(chan amqp.Confirmation, 1))
		for conf := range confirmationCh {
			go p.notifyPublishHandler(Confirmation{
				Confirmation: conf,
			})
		}
	}()
}

func (p *Publisher) startup() error {
	err := declareExchange(p.connection, p.options.ExchangeOptions)
	if err != nil {
		return fmt.Errorf("declare exchange failed: %w", err)
	}
	go p.startNotifyFlowHandler()
	go p.startNotifyBlockedHandler()
	return nil
}

func (p *Publisher) startNotifyFlowHandler() {
	notifyFlowChan := p.connection.ch.NotifyFlow(make(chan bool))
	p.disablePublishDueToFlowMux.Lock()
	p.disablePublishDueToFlow = false
	p.disablePublishDueToFlowMux.Unlock()

	for ok := range notifyFlowChan {
		p.disablePublishDueToFlowMux.Lock()
		if ok {
			logger.Warnf("pausing publishing due to flow request from server")
			p.disablePublishDueToFlow = true
		} else {
			p.disablePublishDueToFlow = false
			logger.Warnf("resuming publishing due to flow request from server")
		}
		p.disablePublishDueToFlowMux.Unlock()
	}
}

func (p *Publisher) startNotifyBlockedHandler() {
	notifyBlockedChan := p.connection.conn.NotifyBlocked(make(chan amqp.Blocking))
	p.disablePublishDueToBlockedMux.Lock()
	p.disablePublishDueToBlocked = false
	p.disablePublishDueToBlockedMux.Unlock()

	for blocked := range notifyBlockedChan {
		p.disablePublishDueToBlockedMux.Lock()
		if blocked.Active {
			logger.Warnf("pausing publishing due to TCP blocking from server")
			p.disablePublishDueToBlocked = true
		} else {
			p.disablePublishDueToBlocked = false
			logger.Warnf("resuming publishing due to TCP blocking from server")
		}
		p.disablePublishDueToBlockedMux.Unlock()
	}
}

func (p *Publisher) Publish(
	data []byte,
	routingKeys []string,
	optionFuncs ...func(*PublishOptions),
) error {
	return p.PublishWithContext(context.Background(), data, routingKeys, optionFuncs...)
}

func (p *Publisher) PublishWithContext(
	ctx context.Context,
	data []byte,
	routingKeys []string,
	optionFuncs ...func(*PublishOptions),
) error {
	p.disablePublishDueToFlowMux.RLock()
	defer p.disablePublishDueToFlowMux.RUnlock()
	if p.disablePublishDueToFlow {
		return fmt.Errorf("publishing blocked due to high flow on the server")
	}

	p.disablePublishDueToBlockedMux.RLock()
	defer p.disablePublishDueToBlockedMux.RUnlock()
	if p.disablePublishDueToBlocked {
		return fmt.Errorf("publishing blocked due to TCP block on the server")
	}

	options := &PublishOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if options.DeliveryMode == 0 {
		options.DeliveryMode = Transient
	}

	message := amqp.Publishing{
		ContentType:     options.ContentType,
		DeliveryMode:    options.DeliveryMode,
		Body:            data,
		Headers:         tableToAMQPTable(options.Headers),
		Expiration:      options.Expiration,
		ContentEncoding: options.ContentEncoding,
		Priority:        options.Priority,
		CorrelationId:   options.CorrelationID,
		ReplyTo:         options.ReplyTo,
		MessageId:       options.MessageID,
		Timestamp:       options.Timestamp,
		Type:            options.Type,
		UserId:          options.UserID,
		AppId:           options.AppID,
	}

	if len(routingKeys) == 0 {
		routingKeys = []string{""}
	}

	for _, routingKey := range routingKeys {
		err := p.connection.ch.PublishWithContext(
			ctx,
			options.Exchange,
			routingKey,
			options.Mandatory,
			options.Immediate,
			message,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Publisher) PublishWithDeferredConfirmWithContext(
	ctx context.Context,
	data []byte,
	routingKeys []string,
	optionFuncs ...func(*PublishOptions),
) (PublisherConfirmation, error) {
	p.disablePublishDueToFlowMux.RLock()
	defer p.disablePublishDueToFlowMux.RUnlock()
	if p.disablePublishDueToFlow {
		return nil, fmt.Errorf("publishing blocked due to high flow on the server")
	}

	p.disablePublishDueToBlockedMux.RLock()
	defer p.disablePublishDueToBlockedMux.RUnlock()
	if p.disablePublishDueToBlocked {
		return nil, fmt.Errorf("publishing blocked due to TCP block on the server")
	}

	options := &PublishOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if options.DeliveryMode == 0 {
		options.DeliveryMode = Transient
	}

	var deferredConfirmations []*amqp.DeferredConfirmation

	if len(routingKeys) == 0 {
		routingKeys = []string{""}
	}

	for _, routingKey := range routingKeys {
		message := amqp.Publishing{}
		message.ContentType = options.ContentType
		message.DeliveryMode = options.DeliveryMode
		message.Body = data
		message.Headers = tableToAMQPTable(options.Headers)
		message.Expiration = options.Expiration
		message.ContentEncoding = options.ContentEncoding
		message.Priority = options.Priority
		message.CorrelationId = options.CorrelationID
		message.ReplyTo = options.ReplyTo
		message.MessageId = options.MessageID
		message.Timestamp = options.Timestamp
		message.Type = options.Type
		message.UserId = options.UserID
		message.AppId = options.AppID

		// Actual publish.
		conf, err := p.connection.ch.PublishWithDeferredConfirmWithContext(
			ctx,
			options.Exchange,
			routingKey,
			options.Mandatory,
			options.Immediate,
			message,
		)
		if err != nil {
			return nil, err
		}
		deferredConfirmations = append(deferredConfirmations, conf)
	}
	return deferredConfirmations, nil
}
