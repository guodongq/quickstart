package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// nullMetric is an OpenTelemetry metric nullMetric.
type nullMetric struct{}

// New returns a metric nullMetric.
func newNullMetric() metric.Exporter {
	return &nullMetric{}
}

func (e *nullMetric) Temporality(_ metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func (e *nullMetric) Aggregation(_ metric.InstrumentKind) metric.Aggregation {
	return metric.AggregationSum{}
}

func (e *nullMetric) Export(_ context.Context, _ *metricdata.ResourceMetrics) error {
	return nil
}

func (e *nullMetric) ForceFlush(_ context.Context) error {
	// nullMetric holds no state, nothing to flush.
	return nil
}

func (e *nullMetric) Shutdown(_ context.Context) error {
	return nil
}
