package decorator

import (
	"context"
	"github.com/guodongq/quickstart/pkg/log"
)

func ApplyCommandDecorators[C any](
	handler CommandHandler[C],
	logger log.Logger,
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
