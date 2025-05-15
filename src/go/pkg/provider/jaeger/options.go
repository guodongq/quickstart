package jaeger

type JaegerOptions struct {
	Enabled       bool
	Host          string
	Port          int
	UseZipkinMode bool
}

func getDefaultJaegerOptions() JaegerOptions {
	return JaegerOptions{
		Enabled:       false,
		Host:          "127.0.0.1",
		Port:          6831,
		UseZipkinMode: false,
	}
}

func WithJaegerOptionsEnabled() func(*JaegerOptions) {
	return func(o *JaegerOptions) {
		o.Enabled = true
	}
}

func WithJaegerOptionsHost(host string) func(*JaegerOptions) {
	return func(o *JaegerOptions) {
		o.Host = host
	}
}

func WithJaegerOptionsPort(port int) func(*JaegerOptions) {
	return func(o *JaegerOptions) {
		o.Port = port
	}
}

func WithJaegerOptionsUseZipkinMode() func(*JaegerOptions) {
	return func(o *JaegerOptions) {
		o.UseZipkinMode = true
	}
}
