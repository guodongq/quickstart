package project

import (
	"github.com/guodongq/quickstart/pkg/ddd"
	"github.com/guodongq/quickstart/pkg/types"
)

// Project represents a project entity
type Project struct {
	*ddd.BaseEntity[string]
	Name        string
	Description string

	Meta *types.Meta
}

func NewProject(optionFuncs ...func(*Project)) *Project {
	var project Project
	for _, optionFunc := range optionFuncs {
		optionFunc(&project)
	}
	return &project
}

func WithProjectBaseEntity(baseEntity *ddd.BaseEntity[string]) func(*Project) {
	return func(project *Project) {
		project.BaseEntity = baseEntity
	}
}

func WithProjectName(name string) func(*Project) {
	return func(project *Project) {
		project.Name = name
	}
}

func WithProjectDescription(description string) func(*Project) {
	return func(project *Project) {
		project.Description = description
	}
}

func WithProjectMeta(meta *types.Meta) func(*Project) {
	return func(project *Project) {
		project.Meta = meta
	}
}
