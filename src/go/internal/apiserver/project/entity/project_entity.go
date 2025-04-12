package entity

import (
	"fmt"
	"time"

	ddd "github.com/guodongq/quickstart/pkg/ddd"
)

type Project[ID comparable] struct {
	ddd.BaseEntity[ID]
	Name        string
	Description string
}

func NewProject(id, name, description string) *Project[string] {
	return &Project[string]{
		//BaseEntity:  ddd.BaseEntity{Id: id},
		Name:        name,
		Description: description,
	}
}

type ProjectCreatedEvent struct {
	ddd.BaseEvent
	ProjectID   any
	ProjectName string
}

func NewProjectCreatedEvent(projectID any, projectName string) *ProjectCreatedEvent {
	return &ProjectCreatedEvent{
		//BaseEvent: *ddd.NewBaseDomainEvent(),
		ProjectID:   projectID,
		ProjectName: projectName,
	}
}

type AggregateRoot struct {
	//ddd.BaseAggregateRoot
	//Project *Project
}

func main() {
	project := NewProject("123", "Zone", "test project")

	event := NewProjectCreatedEvent(project.ID(), project.Name)

	aggregateRoot := AggregateRoot{
		//BaseAggregateRoot: ddd.BaseAggregateRoot{},
		//Project:           project,
	}
	_ = aggregateRoot

	//aggregateRoot.AddDomainEvent(event)

	fmt.Printf("Event occurred on: %s\n", event.OccurredAt().Format(time.RFC3339))
	fmt.Printf("Project ID: %s\n", event.ProjectID)
	fmt.Printf("Project Name: %s\n", event.ProjectName)

	//domainEvents := aggregateRoot.GetDomainEvents()
	//fmt.Printf("Domain Events count: %d\n", len(domainEvents))
	//fmt.Printf("Project Entity: %s\n", aggregateRoot.Project)
}
