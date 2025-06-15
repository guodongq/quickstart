package decorator

import (
	"context"
	"github.com/guodongq/quickstart/pkg/log"
)

func ApplyQueryDecorators[Q any, R any](
	handler QueryHandler[Q, R],
	logger log.Logger,
	metricsClient MetricsClient,
) QueryHandler[Q, R] {
	return queryLoggingDecorator[Q, R]{
		base: queryMetricsDecorator[Q, R]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
