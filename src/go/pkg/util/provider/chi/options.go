package chi

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ChiOptions struct {
	Port           int
	HandlerFromMux func(router chi.Router) http.Handler
	Middlewares    []func(http.Handler) http.Handler
}

func getDefaultChiEngineOptions() ChiOptions {
	return ChiOptions{
		Port: 8080,
	}
}

func WithChiOptionsPort(port int) func(*ChiOptions) {
	return func(o *ChiOptions) {
		o.Port = port
	}
}

func WithChiOptionsHandlerFromMux(f func(router chi.Router) http.Handler) func(*ChiOptions) {
	return func(o *ChiOptions) {
		o.HandlerFromMux = f
	}
}

func WithChiOptionsMiddlewares(middleware ...func(http.Handler) http.Handler) func(*ChiOptions) {
	return func(o *ChiOptions) {
		o.Middlewares = append(o.Middlewares, middleware...)
	}
}
