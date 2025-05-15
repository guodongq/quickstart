package log

import (
	"fmt"
	"strings"
)

// Level defines the priority of a log message.
// When a logger is configured with a Level, any log message with a lower
// log Level (smaller by integer comparison) will not be Output.
type Level int

// The levels of logs.
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
)

// String returns a string representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

func parseLevel(lvl string) (level Level, err error) {
	switch strings.ToLower(lvl) {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warning":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "panic":
		level = PanicLevel
	default:
		err = fmt.Errorf("invalid log level: %s", lvl)
	}

	return
}
