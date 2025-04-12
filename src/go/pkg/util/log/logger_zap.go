package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapToLevel = map[zapcore.Level]Level{
	zapcore.DebugLevel: LevelDebug,
	zapcore.InfoLevel:  LevelInfo,
	zapcore.WarnLevel:  LevelWarn,
	zapcore.ErrorLevel: LevelError,
	zapcore.PanicLevel: LevelPanic,
}

type ZapLogger struct {
	logger  *zap.Logger
	options ZapLoggerOptions
}

func NewZapLogger(optionFuncs ...func(*ZapLoggerOptions)) *ZapLogger {
	defaultOptions := getDefaultZapLoggerOptions()
	options := &defaultOptions

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &ZapLogger{
		options: *options,
		logger: zap.New(
			zapcore.NewCore(
				options.Encoder,
				zapcore.Lock(options.Output),
				zap.NewAtomicLevelAt(options.Level),
			),
		),
	}
}

func (p *ZapLogger) Init() {
	defaultLevel := LevelInfo
	if level, exists := zapToLevel[p.logger.Level()]; exists {
		defaultLevel = level
	}

	SetLevel(defaultLevel)
	SetDefaultLogger(p)
}

func (p *ZapLogger) Debug(v ...any) {
	p.logger.Sugar().Debug(v...)
}

func (p *ZapLogger) Info(v ...any) {
	p.logger.Sugar().Info(v...)
}

func (p *ZapLogger) Warn(v ...any) {
	p.logger.Sugar().Warn(v...)
}

func (p *ZapLogger) Error(v ...any) {
	p.logger.Sugar().Error(v...)
}

func (p *ZapLogger) Panic(v ...any) {
	p.logger.Sugar().Panic(v...)
}

func (p *ZapLogger) Debugf(format string, v ...any) {
	p.logger.Sugar().Debugf(format, v...)
}

func (p *ZapLogger) Infof(format string, v ...any) {
	p.logger.Sugar().Infof(format, v...)
}

func (p *ZapLogger) Warnf(format string, v ...any) {
	p.logger.Sugar().Warnf(format, v...)
}

func (p *ZapLogger) Errorf(format string, v ...any) {
	p.logger.Sugar().Errorf(format, v...)
}

func (p *ZapLogger) Panicf(format string, v ...any) {
	p.logger.Sugar().Panicf(format, v...)
}

func (p *ZapLogger) WithField(key string, value any) Logger {
	return p.WithFields(Fields{key: value})
}

func (p *ZapLogger) WithFields(fields Fields) Logger {
	result := make([]zap.Field, 0)
	for k, v := range fields {
		result = append(result, zap.Any(k, v))
	}
	entry := p.logger.Sugar().Desugar().With(result...)
	return &ZapLogger{
		logger: entry,
	}
}

func (p *ZapLogger) WithError(err error) Logger {
	entry := p.logger.Sugar().Desugar().With(zap.Error(err))
	return &ZapLogger{
		logger: entry,
	}
}

type ZapLoggerOptions struct {
	Level   zapcore.Level
	Encoder zapcore.Encoder
	Output  zapcore.WriteSyncer
}

func getDefaultZapLoggerOptions() ZapLoggerOptions {
	return ZapLoggerOptions{
		Level:   zapcore.InfoLevel,
		Encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		Output:  zapcore.Lock(zapcore.AddSync(os.Stdout)),
	}
}

func WithZapLoggerOptionsLevel(level zapcore.Level) func(*ZapLoggerOptions) {
	return func(options *ZapLoggerOptions) {
		options.Level = level
	}
}

func WithZapLoggerOptionsEncoder(encoder zapcore.Encoder) func(*ZapLoggerOptions) {
	return func(options *ZapLoggerOptions) {
		options.Encoder = encoder
	}
}

func WithZapLoggerOptionsOutput(output zapcore.WriteSyncer) func(*ZapLoggerOptions) {
	return func(options *ZapLoggerOptions) {
		options.Output = output
	}
}
