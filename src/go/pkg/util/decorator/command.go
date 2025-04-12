package decorator

import (
	"context"

	logger "github.com/guodongq/quickstart/pkg/util/log"
)

func ApplyCommandDecorators[C any](
	handler CommandHandler[C],
	logger logger.Logger,
	metricsClient MetricsClient,
) CommandHandler[C] {
	return commandLoggingDecorator[C]{
		base: commandMetricsDecorator[C]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}
