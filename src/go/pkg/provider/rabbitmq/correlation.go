package rabbitmq

import (
	"sync"

	"github.com/google/uuid"
)

type Correlation struct {
	sync.Mutex
	channels map[string]chan Delivery
}

func NewCorrelation() *Correlation {
	return &Correlation{
		channels: make(map[string]chan Delivery),
	}
}

func (c *Correlation) Gen() (correlationID string, val chan Delivery) {
	correlationID = uuid.NewString()
	val = make(chan Delivery, 1)

	c.Lock()
	c.channels[correlationID] = val
	c.Unlock()
	return correlationID, val
}

func (c *Correlation) Del(correlationID string) {
	c.Lock()
	defer c.Unlock()

	ch, exists := c.channels[correlationID]
	if exists {
		close(ch)
		delete(c.channels, correlationID)
	}
}

func (c *Correlation) Put(correlationID string, msg Delivery) {
	c.Lock()
	defer c.Unlock()

	ch, exists := c.channels[correlationID]
	if exists {
		ch <- msg
	}
}
