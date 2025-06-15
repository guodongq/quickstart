package events

type EventRegistry interface {
	Register()
}

type defaultEventRegistryManager struct {
	registries []EventRegistry
}

func NewDefaultEventRegistryManager(
	registers ...EventRegistry,
) EventRegistry {
	return &defaultEventRegistryManager{
		registries: registers,
	}
}
func (r *defaultEventRegistryManager) Register() {
	for _, register := range r.registries {
		register.Register()
	}
}
