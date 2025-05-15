package environment

import "github.com/guodongq/quickstart/pkg/idgen"

type Environment struct {
	ID          idgen.Generator
	Name        string
	Description string
}
