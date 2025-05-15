package decorator

import (
	"context"
	"fmt"
	"github.com/guodongq/quickstart/pkg/log"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger log.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	log := d.logger.WithFields(log.Fields{
		"command":      handlerType,
		"command_body": fmt.Sprintf("%#v", cmd),
	})

	log.Debug("Executing command")
	defer func() {
		if err == nil {
			log.Info("Command executed successfully")
		} else {
			log.WithError(err).Error("Failed to execute command")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[Q any, R any] struct {
	base   QueryHandler[Q, R]
	logger log.Logger
}

func (d queryLoggingDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	log := d.logger.WithFields(log.Fields{
		"query":      generateActionName(query),
		"query_body": fmt.Sprintf("%#v", query),
	})

	log.Debug("Executing query")
	defer func() {
		if err == nil {
			log.Info("Query executed successfully")
		} else {
			log.WithError(err).Error("Failed to execute query")
		}
	}()

	return d.base.Handle(ctx, query)
}
