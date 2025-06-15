package ddd

import "sync"

type AggregateRoot[ID comparable] interface {
	Entity[ID]
	Version() uint
	Events() []Event
	ClearEvents()
}

type BaseAggregateRoot[ID comparable] struct {
	BaseEntity[ID]
	events  []Event
	version uint
	mu      sync.RWMutex
}

func NewBaseAggregateRoot[ID comparable](id ID) *BaseAggregateRoot[ID] {
	return &BaseAggregateRoot[ID]{
		BaseEntity: BaseEntity[ID]{
			Id: id,
		},
	}
}

func (ar *BaseAggregateRoot[ID]) Events() []Event {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return ar.events
}

func (ar *BaseAggregateRoot[ID]) ClearEvents() {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.events = nil
}

func (ar *BaseAggregateRoot[ID]) AddEvent(events ...Event) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.events = append(ar.events, events...)
}

func (ar *BaseAggregateRoot[ID]) Version() uint {
	return ar.version
}

func (ar *BaseAggregateRoot[ID]) IncrementVersion() {
	ar.version++
}
