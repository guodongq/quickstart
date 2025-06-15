package logrus

import (
	"io"
	"os"

	"github.com/guodongq/quickstart/pkg/env"

	"github.com/sirupsen/logrus"
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

func (e EnvConfig) Options() func(*LogrusLoggerOptions) {
	level, err := logrus.ParseLevel(e.level)
	if err != nil {
		level = logrus.InfoLevel
	}

	var formatter logrus.Formatter
	switch e.formatter {
	case "text":
		formatter = &logrus.TextFormatter{}
	case "text_clr":
		formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	case "json":
		fallthrough
	default:
		formatter = &logrus.JSONFormatter{}
	}

	var output io.Writer
	switch e.output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		fallthrough
	default:
		output = os.Stderr
	}

	return func(options *LogrusLoggerOptions) {
		options.Level = level
		options.Formatter = formatter
		options.Output = output
	}
}
