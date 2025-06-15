package bus

import (
	"context"
	"fmt"
	logger "github.com/guodongq/quickstart/pkg/log"
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
	subscribers map[reflect.Type][]reflect.Value
}

func New() EventBus {
	return &DefaultEventBus{
		subscribers: make(map[reflect.Type][]reflect.Value),
	}
}

func (d *DefaultEventBus) Subscribe(fn HandlerFunc) {
	d.Lock()
	defer d.Unlock()

	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		logger.Panicf("subscriber must be a function, got %T", fn)
	}

	fnType := fnVal.Type()
	if fnType.NumIn() != 2 || fnType.NumOut() != 1 {
		logger.Panicf("handler must have signature func(ctx context.Context, event Event) error, got %v", fnType)
	}

	// Check if the handler returns an error.
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !fnType.Out(0).Implements(errorType) {
		logger.Panicf("handler must return an error, got %v", fnType.Out(0))
	}

	// Check if the first parameter is context.Context.
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !fnType.In(0).Implements(ctxType) {
		logger.Panicf("first parameter must be context.Context, got %v", fnType.In(0))
	}

	// Check if the second parameter is a pointer to a struct (event).
	msgType := fnType.In(1)
	if msgType.Kind() != reflect.Ptr || msgType.Elem().Kind() != reflect.Struct {
		logger.Panicf("second parameter must be a pointer to a struct (event), got %v", msgType)
	}

	// Store the subscriber.
	if _, exists := d.subscribers[msgType]; !exists {
		d.subscribers[msgType] = make([]reflect.Value, 0)
	}
	d.subscribers[msgType] = append(d.subscribers[msgType], fnVal)
}

func (d *DefaultEventBus) Publish(ctx context.Context, event Event) error {
	d.RLock()
	defer d.RUnlock()

	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	eventType := reflect.TypeOf(event)
	if eventType.Kind() != reflect.Ptr || eventType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("event must be a pointer to a struct, got %T", event)
	}

	listeners, exists := d.subscribers[eventType]
	if !exists {
		logger.Infof("No subscribers found for event %s", eventType)
		return nil
	}

	params := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(event)}
	return d.notifySubscribers(ctx, listeners, params)
}

func (d *DefaultEventBus) notifySubscribers(ctx context.Context, subscribers []reflect.Value, params []reflect.Value) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, v := range subscribers {
		subscriber := v
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				results := subscriber.Call(params)
				if !results[0].IsNil() {
					if err, ok := results[0].Interface().(error); ok {
						return err
					}
					return fmt.Errorf("subscriber returned non-error type: %T", results[0].Interface())
				}
				return nil
			}
		})
	}
	return g.Wait()
}
