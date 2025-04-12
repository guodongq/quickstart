package log

// Panic calls the default logger's Panic method
func Panic(v ...interface{}) {
	defaultLogger.Panic(v)
}

// Error calls the default logger's Error method.
func Error(v ...interface{}) {
	if level > LevelError {
		return
	}
	defaultLogger.Error(v)
}

// Warn calls the default logger's Warn method.
func Warn(v ...interface{}) {
	if level > LevelWarn {
		return
	}
	defaultLogger.Warn(v)
}

// Info calls the default logger's Info method.
func Info(v ...interface{}) {
	if level > LevelInfo {
		return
	}
	defaultLogger.Info(v)
}

// Debug calls the default logger's Debug method.
func Debug(v ...interface{}) {
	if level > LevelDebug {
		return
	}
	defaultLogger.Debug(v)
}

// Panicf calls the default logger's Fatalf method
func Panicf(format string, v ...interface{}) {
	defaultLogger.Panicf(format, v...)
}

// Errorf calls the default logger's Errorf method.
func Errorf(format string, v ...interface{}) {
	if level > LevelError {
		return
	}
	defaultLogger.Errorf(format, v...)
}

// Warnf calls the default logger's Warnf method.
func Warnf(format string, v ...interface{}) {
	if level > LevelWarn {
		return
	}
	defaultLogger.Warnf(format, v...)
}

// Infof calls the default logger's Infof method.
func Infof(format string, v ...interface{}) {
	if level > LevelInfo {
		return
	}
	defaultLogger.Infof(format, v...)
}

// Debugf calls the default logger's Debugf method.
func Debugf(format string, v ...interface{}) {
	if level > LevelDebug {
		return
	}
	defaultLogger.Debugf(format, v...)
}

// WithField calls the default logger's WithField method.
func WithField(key string, value any) Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields calls the default logger's WithFields method.
func WithFields(fields Fields) Logger {
	return defaultLogger.WithFields(fields)
}

// WithError calls the default logger's WithError method.
func WithError(err error) Logger {
	return defaultLogger.WithError(err)
}
