package http

import "net/http"

type HttpServerOptions struct {
	Port             int
	RegisterHandlers func(*http.ServeMux) http.Handler
}

func getDefaultHttpServerOptions() HttpServerOptions {
	return HttpServerOptions{
		Port: 8080,
	}
}

func WithHttpServerOptionsPort(port int) func(*HttpServerOptions) {
	return func(o *HttpServerOptions) {
		o.Port = port
	}
}

func WithHttpServerOptionsRegisterHandlers(f func(*http.ServeMux) http.Handler) func(*HttpServerOptions) {
	return func(o *HttpServerOptions) {
		o.RegisterHandlers = f
	}
}
