package opentelemetry

import "go.opentelemetry.io/otel/attribute"

type OpenTelemetryOptions struct {
	Enabled              bool
	Endpoint             string
	ServiceName          string
	Attributes           map[string]string
	OpentracingSupported bool
}

func getDefaultOpenTelemetryOptions() OpenTelemetryOptions {
	return OpenTelemetryOptions{
		Enabled:              false,
		Endpoint:             "http://localhost:4317",
		ServiceName:          "app",
		Attributes:           nil,
		OpentracingSupported: false,
	}
}

func WithOpenTelemetryOptionsEnabled() func(*OpenTelemetryOptions) {
	return func(o *OpenTelemetryOptions) {
		o.Enabled = true
	}
}

func WithOpenTelemetryOptionsEndpoint(endpoint string) func(*OpenTelemetryOptions) {
	return func(o *OpenTelemetryOptions) {
		o.Endpoint = endpoint
	}
}

func WithOpenTelemetryOptionsServiceName(serviceName string) func(*OpenTelemetryOptions) {
	return func(o *OpenTelemetryOptions) {
		o.ServiceName = serviceName
	}
}

func WithOpenTelemetryOptionsAttributes(attributes map[string]string) func(*OpenTelemetryOptions) {
	return func(o *OpenTelemetryOptions) {
		o.Attributes = attributes
	}
}

func WithOpenTelemetryOptionsOpentracingSupported() func(*OpenTelemetryOptions) {
	return func(o *OpenTelemetryOptions) {
		o.OpentracingSupported = true
	}
}

func (o *OpenTelemetryOptions) attributeKeyValues() []attribute.KeyValue {
	kvs := make([]attribute.KeyValue, 0)
	for k, v := range o.Attributes {
		kvs = append(kvs, attribute.String(k, v))
	}
	return kvs
}

func (o *OpenTelemetryOptions) MergeIn(opts ...func(*OpenTelemetryOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}
