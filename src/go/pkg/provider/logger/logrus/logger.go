package logrus

import (
	"github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var logrusToLevel = map[logrus.Level]log.Level{
	logrus.DebugLevel: log.DebugLevel,
	logrus.InfoLevel:  log.InfoLevel,
	logrus.WarnLevel:  log.WarnLevel,
	logrus.ErrorLevel: log.ErrorLevel,
	logrus.PanicLevel: log.PanicLevel,
}

type Logrus struct {
	provider.AbstractProvider

	entry   *logrus.Entry
	options LogrusLoggerOptions
}

func New(optionFuncs ...func(*LogrusLoggerOptions)) *Logrus {
	defaultOptions := getDefaultLogrusLoggerOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	LoadEnvConfig().Options()(options)

	logger := logrus.New()
	logger.Out = options.Output
	logger.Formatter = options.Formatter
	logger.Level = options.Level

	return &Logrus{
		entry:   logrus.NewEntry(logger),
		options: *options,
	}
}

func (l *Logrus) Init() error {
	defaultLevel := log.InfoLevel
	if level, exists := logrusToLevel[l.entry.Logger.Level]; exists {
		defaultLevel = level
	}
	log.SetLevel(defaultLevel)
	log.SetDefaultLogger(l)
	return nil
}

func (l *Logrus) Debug(v ...any) {
	l.entry.Debug(v...)
}

func (l *Logrus) Info(v ...any) {
	l.entry.Info(v...)
}

func (l *Logrus) Warn(v ...any) {
	l.entry.Warn(v...)
}

func (l *Logrus) Error(v ...any) {
	l.entry.Error(v...)
}

func (l *Logrus) Panic(v ...any) {
	l.entry.Panic(v...)
}

func (l *Logrus) Debugf(format string, v ...any) {
	l.entry.Debugf(format, v...)
}

func (l *Logrus) Infof(format string, v ...any) {
	l.entry.Infof(format, v...)
}

func (l *Logrus) Warnf(format string, v ...any) {
	l.entry.Warnf(format, v...)
}

func (l *Logrus) Errorf(format string, v ...any) {
	l.entry.Errorf(format, v...)
}

func (l *Logrus) Panicf(format string, v ...any) {
	l.entry.Panicf(format, v...)
}

func (l *Logrus) WithField(key string, value any) log.Logger {
	return &Logrus{
		entry: l.entry.WithField(key, value),
	}
}

func (l *Logrus) WithFields(fields log.Fields) log.Logger {
	logrusFields := make(logrus.Fields, len(fields))
	for k, v := range fields {
		logrusFields[k] = v
	}
	return &Logrus{
		entry: l.entry.WithFields(logrusFields),
	}
}

func (l *Logrus) WithError(err error) log.Logger {
	return &Logrus{
		entry: l.entry.WithError(err),
	}
}

type LogrusLoggerOptions struct {
	Level     logrus.Level
	Formatter logrus.Formatter
	Output    io.Writer
}

func getDefaultLogrusLoggerOptions() LogrusLoggerOptions {
	return LogrusLoggerOptions{
		Level:     logrus.InfoLevel,
		Formatter: &logrus.JSONFormatter{},
		Output:    os.Stderr,
	}
}

func WithLogrusLoggerOptionsLevel(level logrus.Level) func(*LogrusLoggerOptions) {
	return func(options *LogrusLoggerOptions) {
		options.Level = level
	}
}

func WithLogrusLoggerOptionsFormatter(formatter logrus.Formatter) func(*LogrusLoggerOptions) {
	return func(options *LogrusLoggerOptions) {
		options.Formatter = formatter
	}
}

func WithLogrusLoggerOptionsTextFormatter() func(*LogrusLoggerOptions) {
	return func(options *LogrusLoggerOptions) {
		options.Formatter = &logrus.TextFormatter{}
	}
}

func WithLogrusLoggerOptionsTextClrFormatter() func(*LogrusLoggerOptions) {
	return func(options *LogrusLoggerOptions) {
		options.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	}
}

func WithLogrusLoggerOptionsOutput(output io.Writer) func(*LogrusLoggerOptions) {
	return func(options *LogrusLoggerOptions) {
		options.Output = output
	}
}

func WithLogrusLoggerOptionsStdoutOutput() func(*LogrusLoggerOptions) {
	return func(options *LogrusLoggerOptions) {
		options.Output = os.Stdout
	}
}
