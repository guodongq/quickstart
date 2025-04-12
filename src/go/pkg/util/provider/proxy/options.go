package proxy

type ProxyOptions struct {
	Enabled   bool
	Debug     bool
	Port      int
	Endpoint  string
	TargetURL string
	Prefix    string
}

func getDefaultProxyOptions() ProxyOptions {
	return ProxyOptions{
		Enabled:   false,
		Debug:     false,
		Port:      4040,
		Endpoint:  "/",
		TargetURL: "http://localhost:8080",
	}
}

func (o *ProxyOptions) MergeIn(opts ...func(*ProxyOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithProxyOptionsPort(port int) func(*ProxyOptions) {
	return func(o *ProxyOptions) {
		o.Port = port
	}
}

func WithProxyOptionsEnabled() func(*ProxyOptions) {
	return func(o *ProxyOptions) {
		o.Enabled = true
	}
}

func WithProxyOptionsEndpoint(endpoint string) func(*ProxyOptions) {
	return func(o *ProxyOptions) {
		o.Endpoint = endpoint
	}
}

func WithProxyOptionsTargetURL(targetURL string) func(*ProxyOptions) {
	return func(o *ProxyOptions) {
		o.TargetURL = targetURL
	}
}

func WitProxyOptionsDebug() func(*ProxyOptions) {
	return func(o *ProxyOptions) {
		o.Debug = true
	}
}
