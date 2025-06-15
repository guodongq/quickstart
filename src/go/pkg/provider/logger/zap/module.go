package zap

import (
	"github.com/guodongq/quickstart/config"
	"github.com/guodongq/quickstart/pkg/log"
	adapter "github.com/guodongq/quickstart/pkg/provider/logger/fxadapter"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Module = fx.Module(
	"logger",
	fx.WithLogger(func(logger log.Logger) fxevent.Logger {
		return adapter.NewFxEventLogger(logger)
	}),
	fx.Provide(
		func(cfg *config.Config) log.Logger {
			logCfg := cfg.Logger
			level, err := zapcore.ParseLevel(logCfg.Level)
			if err != nil {
				level = zapcore.InfoLevel
			}

			var formatter zapcore.Encoder
			encoder := zap.NewDevelopmentEncoderConfig()
			switch logCfg.Formatter {
			case "json":
				formatter = zapcore.NewJSONEncoder(encoder)
			case "text":
				fallthrough
			default:
				formatter = zapcore.NewConsoleEncoder(encoder)
			}

			var output zapcore.WriteSyncer
			switch logCfg.Output {
			case "stdout":
				output = os.Stdout
			case "stderr":
				fallthrough
			default:
				output = os.Stderr
			}
			return New(
				WithZapLoggerOptionsEncoder(formatter),
				WithZapLoggerOptionsLevel(level),
				WithZapLoggerOptionsOutput(output),
			)
		},
	),
	fx.Invoke(func(logger log.Logger) error {
		type loggerInitializer interface {
			Init() error
		}

		if li, ok := logger.(loggerInitializer); ok {
			return li.Init()
		}
		return nil
	}),
)
