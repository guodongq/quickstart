package log

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

type Fields map[string]any

// Level defines the priority of a log message.
// When a logger is configured with a Level, any log message with a lower
// log Level (smaller by integer comparison) will not be Output.
type Level int

// The levels of logs.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
)

var level Level

// SetLevel sets the Level of logs below which logs will not be Output.
// The default log Level is LevelTrace.
func SetLevel(lv Level) {
	if lv < LevelDebug || lv > LevelPanic {
		panic("invalid Level")
	}
	level = lv
}
