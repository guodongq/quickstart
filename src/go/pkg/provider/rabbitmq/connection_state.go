package rabbitmq

import "sync"

type ConnectionState struct {
	sync.RWMutex
	alive bool
}

func DefaultConnectionState() *ConnectionState {
	return &ConnectionState{
		alive: true,
	}
}

func (s *ConnectionState) IsAlive() bool {
	s.RLock()
	defer s.RUnlock()
	return s.alive
}

func (s *ConnectionState) SetAlive(alive bool) {
	s.Lock()
	defer s.Unlock()
	s.alive = alive
}
