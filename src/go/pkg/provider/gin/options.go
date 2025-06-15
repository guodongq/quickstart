package gin

import "github.com/gin-gonic/gin"

type GinOptions struct {
	Port              int
	RegisterHandlers  func(*gin.Engine) error
	GinMiddlewares    []gin.HandlerFunc
	GinOptionFunc     []gin.OptionFunc
	GinTrustedProxies []string
}

func getDefaultGinOptions() GinOptions {
	return GinOptions{
		Port:             8080,
		RegisterHandlers: nil,
		GinMiddlewares: []gin.HandlerFunc{
			gin.Recovery(),
		},
		GinOptionFunc: []gin.OptionFunc{},
	}
}

func WithGinOptionsPort(port int) func(*GinOptions) {
	return func(o *GinOptions) {
		o.Port = port
	}
}

func WithGinOptionsRegisterHandlers(f func(*gin.Engine) error) func(*GinOptions) {
	return func(o *GinOptions) {
		o.RegisterHandlers = f
	}
}

func WithGinOptionsGinMiddlewares(middlewares ...gin.HandlerFunc) func(*GinOptions) {
	return func(o *GinOptions) {
		o.GinMiddlewares = append(o.GinMiddlewares, middlewares...)
	}
}

func WithGinOptionsGinOptionFunc(optionFunc ...gin.OptionFunc) func(*GinOptions) {
	return func(o *GinOptions) {
		o.GinOptionFunc = append(o.GinOptionFunc, optionFunc...)
	}
}

// WithGinOptionsGinTrustedProxies https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
func WithGinOptionsGinTrustedProxies(proxies ...string) func(*GinOptions) {
	return func(o *GinOptions) {
		o.GinTrustedProxies = append(o.GinTrustedProxies, proxies...)
	}
}

func WithGinOptionsReleaseModeEnabled() func(*GinOptions) {
	return func(o *GinOptions) {
		gin.SetMode(gin.ReleaseMode)
	}
}

func WithGinOptionsDebugModeEnabled() func(*GinOptions) {
	return func(o *GinOptions) {
		gin.SetMode(gin.DebugMode)
	}
}

func WithGinOptionsTestModeEnabled() func(*GinOptions) {
	return func(o *GinOptions) {
		gin.SetMode(gin.TestMode)
	}
}
