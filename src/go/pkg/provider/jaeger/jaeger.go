package jaeger

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"

	"github.com/guodongq/quickstart/pkg/provider/app"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

type Jaeger struct {
	provider.AbstractProvider

	appProvider *app.App
	options     JaegerOptions
	closer      io.Closer
}

func New(appProvider *app.App, optionFuncs ...func(*JaegerOptions)) *Jaeger {
	defaultOptions := getDefaultJaegerOptions()
	options := &defaultOptions

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	return &Jaeger{
		appProvider: appProvider,
		options:     *options,
	}
}

func (p *Jaeger) Init() error {
	if p.appProvider == nil {
		return fmt.Errorf("app provider is required")
	}

	metrics := prometheus.New()
	conf, err := config.FromEnv()
	if err != nil {
		logger.WithError(err).Error("Jaeger Config From Env launch failed")
		return err
	}
	conf.ServiceName = p.appProvider.Name()
	conf.Disabled = !p.options.Enabled
	conf.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	conf.Reporter = &config.ReporterConfig{
		LogSpans:           true,
		LocalAgentHostPort: fmt.Sprintf("%s:%d", p.options.Host, p.options.Port),
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator(
		zipkin.BaggagePrefix("x-request-id"),
	)
	opts := []config.Option{
		config.Metrics(metrics),
		config.Logger(newJaegerLogger(logger.DefaultLogger())),
		config.ZipkinSharedRPCSpan(true),
	}

	if p.options.UseZipkinMode {
		opts = append(opts,
			config.Injector(opentracing.HTTPHeaders, zipkinPropagator),
			config.Extractor(opentracing.HTTPHeaders, zipkinPropagator))
	}

	// Use the configuration to create a new tracer.
	tracer, closer, err := conf.NewTracer(opts...)
	if err != nil {
		logger.WithError(err).Error("Jaeger Tracer Provider launch failed")
		return nil
	}

	opentracing.SetGlobalTracer(tracer)
	p.closer = closer

	return nil
}

// Close closes the connection to Jaeger, using the closer created during startup.
func (p *Jaeger) Close() error {
	err := p.closer.Close()
	if err != nil {
		logger.WithError(err).Info("Error while closing Jaeger tracer")
		return err
	}

	return p.AbstractProvider.Close()
}

type jaegerLogger struct {
	logger logger.Logger
}

func newJaegerLogger(logger logger.Logger) *jaegerLogger {
	return &jaegerLogger{logger: logger}
}

// Error logs a message at error priority
func (j *jaegerLogger) Error(msg string) {
	j.logger.Error(msg)
}

// Infof logs a message at info priority
func (j *jaegerLogger) Infof(msg string, args ...interface{}) {
	j.logger.Infof(msg, args...)
}
