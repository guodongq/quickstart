package app

import (
	"path"

	"github.com/guodongq/quickstart/internal/common"
	logger "github.com/guodongq/quickstart/pkg/util/log"
	"github.com/guodongq/quickstart/pkg/util/provider"
)

type App struct {
	provider.AbstractProvider

	options AppOptions
}

func New(optionFuncs ...func(*AppOptions)) *App {
	defaultOptions := getDefaultAppOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &App{
		options: *options,
	}
}

func (p *App) Init() error {
	logger.WithFields(logger.Fields{
		"name":    p.Name(),
		"version": p.Version().String(),
	}).Info("App Provider initialized")
	return nil
}

func (p *App) Name() string {
	return p.options.Name
}

func (p *App) Version() common.Version {
	return common.CurrentVersion()
}

func (p *App) ParseEndpoint(elem ...string) string {
	elem = append([]string{p.options.BasePath}, elem...)
	return path.Join(elem...)
}

func (p *App) ParsePath(elem ...string) string {
	res := p.ParseEndpoint(elem...)
	if res != "/" {
		res += "/"
	}
	return res
}
