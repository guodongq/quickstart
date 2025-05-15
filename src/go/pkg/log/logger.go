package log

type Fields map[string]any

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	Debug(v ...any)
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Panic(v ...any)

	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	Panicf(format string, v ...any)

	WithField(key string, value any) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
}
