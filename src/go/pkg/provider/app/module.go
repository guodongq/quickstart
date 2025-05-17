package app

import (
	"github.com/guodongq/quickstart/config"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	fx.Provide(func(cfg config.Config) *App {
		return New(
			WithAppOptionsName(cfg.Server.Name),
			WithAppOptionsBasePath(cfg.Server.BasePath),
		)
	}),
	fx.Invoke(func(app *App) error {
		return app.Init()
	}),
)
