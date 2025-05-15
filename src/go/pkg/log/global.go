package log

import "runtime/debug"

// Panic calls the default logger's Panic method
func Panic(v ...interface{}) {
	if shouldLog(PanicLevel) {
		DefaultLogger().Panic(v...)
		debug.PrintStack()
	}
}

// Error calls the default logger's Error method.
func Error(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		DefaultLogger().Error(v...)
	}
}

// Warn calls the default logger's Warn method.
func Warn(v ...interface{}) {
	if shouldLog(WarnLevel) {
		DefaultLogger().Warn(v...)
	}
}

// Info calls the default logger's Info method.
func Info(v ...interface{}) {
	if shouldLog(InfoLevel) {
		DefaultLogger().Info(v...)
	}
}

// Debug calls the default logger's Debug method.
func Debug(v ...interface{}) {
	if shouldLog(DebugLevel) {
		DefaultLogger().Debug(v...)
	}
}

// Panicf calls the default logger's Fatalf method
func Panicf(format string, v ...interface{}) {
	if shouldLog(PanicLevel) {
		DefaultLogger().Panicf(format, v...)
	}
}

// Errorf calls the default logger's Errorf method.
func Errorf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		DefaultLogger().Errorf(format, v...)
	}
}

// Warnf calls the default logger's Warnf method.
func Warnf(format string, v ...interface{}) {
	if shouldLog(WarnLevel) {
		DefaultLogger().Warnf(format, v...)
	}
}

// Infof calls the default logger's Infof method.
func Infof(format string, v ...interface{}) {
	if shouldLog(InfoLevel) {
		DefaultLogger().Infof(format, v...)
	}
}

// Debugf calls the default logger's Debugf method.
func Debugf(format string, v ...interface{}) {
	if shouldLog(DebugLevel) {
		DefaultLogger().Debugf(format, v...)
	}
}

// WithField calls the default logger's WithField method.
func WithField(key string, value any) Logger {
	return DefaultLogger().WithField(key, value)
}

// WithFields calls the default logger's WithFields method.
func WithFields(fields Fields) Logger {
	return DefaultLogger().WithFields(fields)
}

// WithError calls the default logger's WithError method.
func WithError(err error) Logger {
	return DefaultLogger().WithError(err)
}

// shouldLog checks if the given level should be logged based on the global level.
func shouldLog(lvl Level) bool {
	mu.RLock()
	defer mu.RUnlock()
	return lvl >= level
}
