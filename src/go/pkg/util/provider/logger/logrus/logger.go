package logrus

import (
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
)

type LogrusWrapper struct {
	provider.AbstractProvider

	logger *logger.LogrusLogger
}

func New(optionFuncs ...func(*logger.LogrusLoggerOptions)) *LogrusWrapper {
	logrusLogger := logger.NewLogrusLogger(optionFuncs...)
	return &LogrusWrapper{
		logger: logrusLogger,
	}
}

func (l *LogrusWrapper) Init() error {
	l.logger.Init()
	return nil
}
