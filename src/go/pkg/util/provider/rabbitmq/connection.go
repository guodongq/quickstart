package rabbitmq

import (
	"fmt"

	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
	"github.com/guodongq/quickstart/pkg/util/provider/probes"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	provider.AbstractProvider
	probesProvider *probes.Probes
	options        ConnectionOptions

	conn      *amqp.Connection
	connState *ConnectionState

	ch      *amqp.Channel
	chState *ConnectionState
}

type Config amqp.Config

func NewConnection(probesProvider *probes.Probes, optionFuncs ...func(*ConnectionOptions)) *Connection {
	defaultOptions := getDefaultConnectionOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Connection{
		probesProvider: probesProvider,
		options:        *options,
		connState:      DefaultConnectionState(),
		chState:        DefaultConnectionState(),
	}
}

func (c *Connection) Init() (err error) {
	if len(c.options.URL) == 0 {
		return fmt.Errorf("rabbitmq url is required")
	}

	c.conn, err = amqp.DialConfig(c.options.URL, amqp.Config(c.options.Config))
	if err != nil {
		return err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		_ = c.conn.Close()
		return err
	}

	go c.notifyChanClose()
	go c.notifyConnClose()

	if c.probesProvider != nil {
		c.probesProvider.AddLivenessProbes(c.livenessProbe)
	}

	return nil
}

func (c *Connection) Close() error {
	if c.ch != nil {

		_ = c.ch.Close()
	}

	if c.conn != nil {
		_ = c.conn.Close()
	}

	return nil
}

func (c *Connection) livenessProbe() error {
	if c.chState.IsAlive() && c.connState.IsAlive() {
		return nil
	}

	return amqp.ErrClosed
}

func (c *Connection) notifyConnClose() {
	notifyCloseChan := c.conn.NotifyClose(make(chan *amqp.Error, 1))

	err := <-notifyCloseChan

	c.connState.SetAlive(false)
	if err != nil {
		logger.WithError(err).Errorf("amqp: connection closed")
	}
	if err == nil {
		logger.Warnf("amqp connection closed gracefully")
	}
}

func (c *Connection) notifyChanClose() {
	notifyCloseChan := c.ch.NotifyClose(make(chan *amqp.Error, 1))

	err := <-notifyCloseChan

	c.chState.SetAlive(false)
	if err != nil {
		logger.WithError(err).Errorf("amqp: channel closed")
	}
	if err == nil {
		logger.Warnf("amqp channel closed gracefully")
	}
}
