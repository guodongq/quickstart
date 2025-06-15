package environment

import "context"

type Repository interface {
	CreateEnvironment(ctx context.Context, entity *Environment) error
}
