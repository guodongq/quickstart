package prometheus

type PrometheusOptions struct {
	Port     int
	Enabled  bool
	Endpoint string
}

func getDefaultPrometheusOptions() PrometheusOptions {
	return PrometheusOptions{
		Port:     9090,
		Enabled:  false,
		Endpoint: "/metrics",
	}
}

func (o *PrometheusOptions) MergeIn(opts ...func(*PrometheusOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithPrometheusOptionsPort(port int) func(*PrometheusOptions) {
	return func(o *PrometheusOptions) {
		o.Port = port
	}
}

func WithPrometheusOptionsEnabled() func(*PrometheusOptions) {
	return func(o *PrometheusOptions) {
		o.Enabled = true
	}
}

func WithPrometheusOptionsEndpoint(endpoint string) func(*PrometheusOptions) {
	return func(o *PrometheusOptions) {
		o.Endpoint = endpoint
	}
}
