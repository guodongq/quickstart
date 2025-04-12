package eventbus

import (
	"context"
	"fmt"
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Event defines a event interface.
type Event any

// HandlerFunc defines a handler function interface.
type HandlerFunc any

// EventBus type defines the bus interface structure.
type EventBus interface {
	Subscribe(fn HandlerFunc)
	Publish(ctx context.Context, event Event) error
}

type DefaultEventBus struct {
	sync.RWMutex
	subscribers map[string][]HandlerFunc
}

func New() EventBus {
	return &DefaultEventBus{
		subscribers: make(map[string][]HandlerFunc),
	}
}

func (d *DefaultEventBus) Subscribe(fn HandlerFunc) {
	d.Lock()
	defer d.Unlock()
	handlerType := reflect.TypeOf(fn)
	// check handler takes two args and returns zero error
	if handlerType.NumIn() != 2 {
		logger.Panicf("expected handler <%v> to take two args, a context and a event structure pointer", fn)
	}
	if handlerType.NumOut() != 1 {
		logger.Panicf("expected handler <%v> to return one arg, an error", fn)
	}

	// check handler returns an error
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !handlerType.Out(0).Implements(errorType) {
		logger.Panicf("expected handler <%v> to return one arg, an error", fn)
	}

	// check first arg is a context.Context
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !handlerType.In(0).Implements(ctxType) {
		logger.Panicf("expected handler <%v> to take two args, a context and a msg", fn)
	}

	// check second arg is a struct pointer
	msgType := handlerType.In(1)
	if msgType.Kind() != reflect.Ptr || msgType.Elem().Kind() != reflect.Struct {
		logger.Panicf("expected handler <%v> to take two args, a context and a event structure pointer", fn)
	}

	topic := msgType.Elem().Name()
	_, exists := d.subscribers[topic]
	if !exists {
		d.subscribers[topic] = make([]HandlerFunc, 0)
	}
	d.subscribers[topic] = append(d.subscribers[topic], fn)
}

func (d *DefaultEventBus) Publish(ctx context.Context, event Event) error {
	d.RLock()
	defer d.RUnlock()
	eventPtr := reflect.TypeOf(event)
	if eventPtr.Kind() != reflect.Ptr || eventPtr.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected published event to be a struct pointer, got %T", event)
	}

	topic := eventPtr.Elem().Name()

	var params []reflect.Value
	if listeners, exists := d.subscribers[topic]; exists {
		params = append(params, reflect.ValueOf(ctx))
		params = append(params, reflect.ValueOf(event))
		if err := d.notifySubscribers(listeners, params); err != nil {
			return err
		}
	}
	return nil
}

func (d *DefaultEventBus) notifySubscribers(subscribers []HandlerFunc, params []reflect.Value) error {
	g, _ := errgroup.WithContext(context.Background())
	for _, v := range subscribers {
		subscriber := v
		g.Go(func() error {
			result := reflect.ValueOf(subscriber).Call(params)
			e := result[0].Interface()
			if e != nil {
				err, ok := e.(error)
				if ok {
					return err
				}
				return fmt.Errorf("expected subscriber to return an error, got '%T'", e)
			}
			return nil
		})
	}
	return g.Wait()
}
