package zap

import (
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"os"

	"github.com/guodongq/quickstart/pkg/util/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerLevelEnvKey = []string{
		"LOGGER_LEVEL",
	}

	loggerFormatterEnvKey = []string{
		"LOGGER_FORMATTER",
	}

	loggerOutputEnvKey = []string{
		"LOGGER_OUTPUT",
	}
)

type EnvConfig struct {
	level     string
	formatter string
	output    string
}

func LoadEnvConfig() EnvConfig {
	cfg := EnvConfig{}
	env.SetFromEnvVal(&cfg.level, loggerLevelEnvKey)
	env.SetFromEnvVal(&cfg.formatter, loggerFormatterEnvKey)
	env.SetFromEnvVal(&cfg.output, loggerOutputEnvKey)

	return cfg
}

func (e EnvConfig) Options() func(*logger.ZapLoggerOptions) {
	level, err := zapcore.ParseLevel(e.level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	var formatter zapcore.Encoder
	encoder := zap.NewProductionEncoderConfig()
	switch e.formatter {
	case "json":
		formatter = zapcore.NewJSONEncoder(encoder)
	case "text":
		fallthrough
	default:
		formatter = zapcore.NewConsoleEncoder(encoder)
	}

	var output zapcore.WriteSyncer
	switch e.output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		fallthrough
	default:
		output = os.Stderr
	}

	return func(options *logger.ZapLoggerOptions) {
		options.Level = level
		options.Encoder = formatter
		options.Output = output
	}
}
