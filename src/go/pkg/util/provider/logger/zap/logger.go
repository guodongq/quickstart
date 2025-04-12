package zap

import (
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
)

type ZapWrapper struct {
	provider.AbstractProvider

	logger *logger.ZapLogger
}

func New(optionFuncs ...func(*logger.ZapLoggerOptions)) *ZapWrapper {
	zapLogger := logger.NewZapLogger(optionFuncs...)
	return &ZapWrapper{
		logger: zapLogger,
	}
}

func (l *ZapWrapper) Init() error {
	l.logger.Init()
	return nil
}
