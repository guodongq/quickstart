package project

import (
	"github.com/guodongq/quickstart/pkg/ddd"
	"github.com/guodongq/quickstart/pkg/idgen"
	"github.com/guodongq/quickstart/pkg/types"
)

// Project represents a project entity
type Project struct {
	ddd.BaseEntity[idgen.Generator]
	Name        string
	Description string

	Limitation Limitation
	Metrics    Metrics

	Meta types.Meta
}

type (
	Limitation struct {
		Cpu    int64
		Memory int64
	}

	Metrics struct {
	}
)
