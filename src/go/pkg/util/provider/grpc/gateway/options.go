package gateway

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type GatewayOptions struct {
	Enabled bool
	Port    int

	IncomingHeaderMatcher runtime.HeaderMatcherFunc
}

func getDefaultGatewayOptions() GatewayOptions {
	return GatewayOptions{
		Enabled: false,
		Port:    8080,
		IncomingHeaderMatcher: func(s string) (string, bool) {
			switch s {
			case
				"X-Tenant-ID",
				"X-User-ID":
				return strings.ToLower(s), true
			}

			return runtime.DefaultHeaderMatcher(s)
		},
	}
}

func WithGatewayOptionsIncomingHeaderMatcher(incomingHeaderMatcher runtime.HeaderMatcherFunc) func(*GatewayOptions) {
	return func(o *GatewayOptions) {
		o.IncomingHeaderMatcher = incomingHeaderMatcher
	}
}

func WithGatewayOptionsPort(port int) func(*GatewayOptions) {
	return func(o *GatewayOptions) {
		o.Port = port
	}
}

func WithGatewayOptionsEnabled() func(*GatewayOptions) {
	return func(o *GatewayOptions) {
		o.Enabled = true
	}
}
