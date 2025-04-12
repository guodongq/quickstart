package decorator

import (
	"context"

	logger "github.com/guodongq/quickstart/pkg/util/log"
)

func ApplyQueryDecorators[Q any, R any](
	handler QueryHandler[Q, R],
	logger logger.Logger,
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
