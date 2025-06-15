package app

import (
	"path"

	logger "github.com/guodongq/quickstart/pkg/log"
	"github.com/guodongq/quickstart/pkg/provider"
	"github.com/guodongq/quickstart/pkg/version"
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

func (p *App) Version() version.Version {
	return version.CurrentVersion()
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
