package project

import "github.com/guodongq/quickstart/pkg/ddd"

type ProjectCreatedEvent struct {
	ddd.BaseEvent
	ProjectID   any
	ProjectName string
}

func (p *ProjectCreatedEvent) EventType() string {
	return "project.created"
}
