package pprof

type PProfOptions struct {
	Port     int
	Enabled  bool
	Endpoint string
}

func getDefaultOptions() PProfOptions {
	return PProfOptions{
		Port:     9999,
		Enabled:  false,
		Endpoint: "/debug/pprof",
	}
}

func WithPProfOptionsPort(port int) func(*PProfOptions) {
	return func(o *PProfOptions) {
		o.Port = port
	}
}

func WithPProfOptionsEnabled() func(*PProfOptions) {
	return func(o *PProfOptions) {
		o.Enabled = true
	}
}

func WithPProfOptionsEndpoint(endpoint string) func(*PProfOptions) {
	return func(o *PProfOptions) {
		o.Endpoint = endpoint
	}
}
