package probes

type ProbesOptions struct {
	Port              int
	Enabled           bool
	LivenessEndpoint  string
	ReadinessEndpoint string
}

func getDefaultProbesOptions() ProbesOptions {
	return ProbesOptions{
		Port:              8000,
		Enabled:           true,
		LivenessEndpoint:  "/healthz",
		ReadinessEndpoint: "/ready",
	}
}

func (o *ProbesOptions) MergeIn(opts ...func(*ProbesOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithProbesOptionsPort(port int) func(*ProbesOptions) {
	return func(o *ProbesOptions) {
		o.Port = port
	}
}

func WithProbesOptionsEnabled() func(*ProbesOptions) {
	return func(o *ProbesOptions) {
		o.Enabled = true
	}
}

func WithProbesOptionsDisabled() func(*ProbesOptions) {
	return func(o *ProbesOptions) {
		o.Enabled = false
	}
}

func WithProbesOptionsEnable(enabled bool) func(*ProbesOptions) {
	return func(o *ProbesOptions) {
		o.Enabled = enabled
	}
}

func WithProbesOptionsLiveEndpoint(endpoint string) func(*ProbesOptions) {
	return func(o *ProbesOptions) {
		o.LivenessEndpoint = endpoint
	}
}

func WithProbesOptionsReadEndpoint(endpoint string) func(*ProbesOptions) {
	return func(o *ProbesOptions) {
		o.ReadinessEndpoint = endpoint
	}
}
