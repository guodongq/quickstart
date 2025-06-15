package gin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/guodongq/quickstart/pkg/provider"
)

type Gin struct {
	provider.AbstractRunProvider

	options GinOptions
	engine  *gin.Engine
}

func New(optionFuncs ...func(*GinOptions)) *Gin {
	defaultOptions := getDefaultGinOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Gin{
		options: *options,
	}
}

func (g *Gin) Init() error {
	g.engine = gin.New(g.options.GinOptionFunc...)
	g.engine.Use(g.options.GinMiddlewares...)
	err := g.engine.SetTrustedProxies(g.options.GinTrustedProxies)
	if err != nil {
		return err
	}

	if g.options.RegisterHandlers != nil {
		return g.options.RegisterHandlers(g.engine)
	}

	return nil
}

func (g *Gin) Run() error {
	addr := fmt.Sprintf("127.0.0.1:%d", g.options.Port)
	g.SetRunning(true)

	if err := g.engine.Run(addr); err != nil {
		return err
	}
	return nil
}
