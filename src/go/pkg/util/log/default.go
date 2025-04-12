package log

var defaultLogger Logger = &nopLogger{}

func DefaultLogger() Logger {
	return defaultLogger
}

// SetDefaultLogger sets the default logger.
// This is not concurrency safe, which means it should only be called during init.
func SetDefaultLogger(l Logger) {
	if l == nil {
		panic("logger must not be nil")
	}
	defaultLogger = l
}

var _ Logger = (*nopLogger)(nil)

type nopLogger struct{}

func (l *nopLogger) Panic(_ ...interface{}) {
}

func (l *nopLogger) Error(_ ...interface{}) {
}

func (l *nopLogger) Warn(_ ...interface{}) {
}

func (l *nopLogger) Info(_ ...interface{}) {
}

func (l *nopLogger) Debug(_ ...interface{}) {
}

func (l *nopLogger) Panicf(_ string, _ ...interface{}) {
}

func (l *nopLogger) Errorf(_ string, _ ...interface{}) {
}

func (l *nopLogger) Warnf(_ string, _ ...interface{}) {
}

func (l *nopLogger) Infof(_ string, _ ...interface{}) {
}

func (l *nopLogger) Debugf(_ string, _ ...interface{}) {
}

func (l *nopLogger) WithField(_ string, _ any) Logger {
	return l
}

func (l *nopLogger) WithFields(_ Fields) Logger {
	return l
}

func (l *nopLogger) WithError(_ error) Logger {
	return l
}
