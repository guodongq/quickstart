package environment

type EnvironmentCreatedEvent struct {
	environment Environment
}

func (p *EnvironmentCreatedEvent) EventType() string {
	return "environment.created"
}
