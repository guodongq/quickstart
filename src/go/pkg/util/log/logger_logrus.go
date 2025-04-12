package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logrusToLevel = map[logrus.Level]Level{
	logrus.DebugLevel: LevelDebug,
	logrus.InfoLevel:  LevelInfo,
	logrus.WarnLevel:  LevelWarn,
	logrus.ErrorLevel: LevelError,
	logrus.PanicLevel: LevelPanic,
}

type LogrusLogger struct {
	entry   *logrus.Entry
	options LogrusLoggerOptions
}

func NewLogrusLogger(optionFuncs ...func(*LogrusLoggerOptions)) *LogrusLogger {
	defaultOptions := getDefaultLogrusLoggerOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &LogrusLogger{
		entry: logrus.NewEntry(&logrus.Logger{
			Out:          options.Output,
			Hooks:        make(logrus.LevelHooks),
			Formatter:    options.Formatter,
			ReportCaller: false,
			Level:        options.Level,
			ExitFunc:     os.Exit,
		}),
		options: *options,
	}
}

func (l *LogrusLogger) Init() {
	defaultLevel := LevelInfo
	if level, exists := logrusToLevel[l.entry.Logger.Level]; exists {
		defaultLevel = level
	}
	SetLevel(defaultLevel)
	SetDefaultLogger(l)
}

func (l *LogrusLogger) Debug(v ...any) {
	l.entry.Debug(v...)
}

func (l *LogrusLogger) Info(v ...any) {
	l.entry.Info(v...)
}

func (l *LogrusLogger) Warn(v ...any) {
	l.entry.Warn(v...)
}

func (l *LogrusLogger) Error(v ...any) {
	l.entry.Error(v...)
}

func (l *LogrusLogger) Panic(v ...any) {
	l.entry.Panic(v...)
}

func (l *LogrusLogger) Debugf(format string, v ...any) {
	l.entry.Debugf(format, v...)
}

func (l *LogrusLogger) Infof(format string, v ...any) {
	l.entry.Infof(format, v...)
}

func (l *LogrusLogger) Warnf(format string, v ...any) {
	l.entry.Warnf(format, v...)
}

func (l *LogrusLogger) Errorf(format string, v ...any) {
	l.entry.Errorf(format, v...)
}

func (l *LogrusLogger) Panicf(format string, v ...any) {
	l.entry.Panicf(format, v...)
}

func (l *LogrusLogger) WithField(key string, value any) Logger {
	return l.WithFields(Fields{key: value})
}

func (l *LogrusLogger) WithFields(fields Fields) Logger {
	result := make(logrus.Fields)
	for k, v := range fields {
		result[k] = v
	}

	entry := l.entry.Dup().WithFields(result)
	return &LogrusLogger{
		entry: entry,
	}
}

func (l *LogrusLogger) WithError(err error) Logger {
	entry := l.entry.Dup().WithError(err)
	return &LogrusLogger{
		entry: entry,
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
