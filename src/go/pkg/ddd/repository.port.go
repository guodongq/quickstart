package ddd

import (
	"context"
	"github.com/guodongq/quickstart/pkg/types"
)

type Repository[ID comparable, AR AggregateRoot[ID]] interface {
	Getter[ID, AR]
	Saver[ID, AR]
	Remover[ID, AR]
	Updater[ID, AR]
	Querier[ID, AR]
	UnitOfWorkHandler
}

type Getter[ID comparable, AR AggregateRoot[ID]] interface {
	Get(ctx context.Context, id ID) (AR, error)
}

type Saver[ID comparable, AR AggregateRoot[ID]] interface {
	Save(ctx context.Context, aggregate AR) error
}

type Remover[ID comparable, AR AggregateRoot[ID]] interface {
	Delete(ctx context.Context, id ID) error
	SoftDelete(ctx context.Context, id ID) error
}

type Updater[ID comparable, AR AggregateRoot[ID]] interface {
	Update(ctx context.Context, id ID, updateFn func(context.Context, AR) (AR, error)) error
}

type Querier[ID comparable, AR AggregateRoot[ID]] interface {
	Find(ctx context.Context, spec types.Specification[AR], options ...types.QueryOption) ([]AR, error)
	FindOne(ctx context.Context, spec types.Specification[AR]) (AR, error)
	Count(ctx context.Context, spec types.Specification[AR]) (int64, error)
}

type UnitOfWorkHandler interface {
	Begin(ctx context.Context) (UnitOfWork, error)
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type UnitOfWork interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	IsActive() bool
}
