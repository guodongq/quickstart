package zap

import (
	"github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var zapToLevel = map[zapcore.Level]log.Level{
	zapcore.DebugLevel: log.DebugLevel,
	zapcore.InfoLevel:  log.InfoLevel,
	zapcore.WarnLevel:  log.WarnLevel,
	zapcore.ErrorLevel: log.ErrorLevel,
	zapcore.PanicLevel: log.PanicLevel,
}

type Zap struct {
	provider.AbstractProvider

	logger  *zap.Logger
	sugared *zap.SugaredLogger
	options ZapLoggerOptions
}

func New(optionFuncs ...func(*ZapLoggerOptions)) *Zap {
	defaultOptions := getDefaultZapLoggerOptions()
	options := &defaultOptions

	LoadEnvConfig().Options()(options)

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	core := zapcore.NewCore(
		options.Encoder,
		options.Output,
		options.Level,
	)
	logger := zap.New(core, zap.AddCallerSkip(1))
	sugared := logger.Sugar()

	return &Zap{
		options: *options,
		logger:  logger,
		sugared: sugared,
	}
}

func (p *Zap) Init() error {
	defaultLevel := log.InfoLevel
	if level, exists := zapToLevel[p.logger.Level()]; exists {
		defaultLevel = level
	}

	log.SetLevel(defaultLevel)
	log.SetDefaultLogger(p)
	return nil
}

func (p *Zap) Debug(v ...any) {
	p.sugared.Debug(v...)
}

func (p *Zap) Info(v ...any) {
	p.sugared.Info(v...)
}

func (p *Zap) Warn(v ...any) {
	p.sugared.Warn(v...)
}

func (p *Zap) Error(v ...any) {
	p.sugared.Error(v...)
}

func (p *Zap) Panic(v ...any) {
	p.sugared.Panic(v...)
}

func (p *Zap) Debugf(format string, v ...any) {
	p.sugared.Debugf(format, v...)
}

func (p *Zap) Infof(format string, v ...any) {
	p.sugared.Infof(format, v...)
}

func (p *Zap) Warnf(format string, v ...any) {
	p.sugared.Warnf(format, v...)
}

func (p *Zap) Errorf(format string, v ...any) {
	p.sugared.Errorf(format, v...)
}

func (p *Zap) Panicf(format string, v ...any) {
	p.sugared.Panicf(format, v...)
}

func (p *Zap) WithField(key string, value any) log.Logger {
	return &Zap{
		logger:  p.logger.With(zap.Any(key, value)),
		sugared: p.sugared.With(key, value),
	}
}

func (p *Zap) WithFields(fields log.Fields) log.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Zap{
		logger:  p.logger.With(zapFields...),
		sugared: p.sugared.With(fields),
	}
}

func (p *Zap) WithError(err error) log.Logger {
	return &Zap{
		logger:  p.logger.With(zap.Error(err)),
		sugared: p.sugared.With("error", err),
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
		Encoder: zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
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
