package project

import "github.com/guodongq/quickstart/pkg/ddd"

type ProjectCreatedEvent struct {
	ddd.BaseEvent
	ProjectID   any
	ProjectName string
}

func NewProjectCreatedEvent(projectID any, projectName string) *ProjectCreatedEvent {
	event := &ProjectCreatedEvent{
		ProjectID:   projectID,
		ProjectName: projectName,
	}
	event.BaseEvent = *ddd.NewBaseEvent(event.EventType())
	return event
}

func (p *ProjectCreatedEvent) EventType() string {
	return "project.created"
}
