package opentelemetry

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
	propagation "go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	bridgeopentracing "go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type OpenTelemetry struct {
	provider.AbstractProvider

	options OpenTelemetryOptions
	close   func() error
	flush   func() error
}

func New(optionFuncs ...func(*OpenTelemetryOptions)) *OpenTelemetry {
	defaultOptions := getDefaultOpenTelemetryOptions()
	options := &defaultOptions

	options.MergeIn(optionFuncs...)

	return &OpenTelemetry{
		options: *options,
	}
}

func (p *OpenTelemetry) Init() error {
	if !p.options.Enabled {
		logger.Info("OpenTelemetry provider is disabled")
		return nil
	}

	exporter, err := otlptrace.New(context.Background(), otlptracegrpc.NewClient())
	if err != nil {
		return err
	}
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			p.options.attributeKeyValues()...,
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	// ignore the metrics
	meterReader := metricsdk.NewPeriodicReader(newNullMetric())

	meterProvider := metricsdk.NewMeterProvider(metricsdk.WithReader(meterReader))
	otel.SetMeterProvider(meterProvider)

	if p.options.OpentracingSupported {
		bridgeTracer, wrapperTracerProvider := bridgeopentracing.NewTracerPair(
			tracerProvider.Tracer(""),
		)
		otel.SetTracerProvider(wrapperTracerProvider)
		opentracing.SetGlobalTracer(bridgeTracer)
	}

	otel.SetTextMapPropagator(propagation.NewTextMapPropagator())

	p.close = func() error {
		_ = meterProvider.Shutdown(context.Background())
		return tracerProvider.Shutdown(context.Background())
	}

	p.flush = func() error {
		_ = meterProvider.ForceFlush(context.Background())
		return tracerProvider.ForceFlush(context.Background())
	}

	return nil
}

func (p *OpenTelemetry) Enabled() bool {
	return p.options.Enabled
}

func (p *OpenTelemetry) Flush() error {
	if p.flush != nil {
		return p.flush()
	}
	return errors.New("no instance")
}

func (p *OpenTelemetry) Close() error {
	if p.Enabled() {
		if p.flush != nil {
			return p.flush()
		}
		if p.close != nil {
			return p.close()
		}
	}
	return nil
}

func DefaultTracer() trace.Tracer {
	return otel.Tracer("")
}
