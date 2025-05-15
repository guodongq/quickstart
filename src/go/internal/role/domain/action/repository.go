package action

import "context"

type Repository interface {
	AddAction(ctx context.Context, action *Action) error
	GetAction(ctx context.Context, actionUUID string) (*Action, error)
	UpdateAction(
		ctx context.Context,
		actionUUID string,
		updateFn func(ctx context.Context, action *Action) (*Action, error),
	) error
}
