package probes

import (
	"context"
	"github.com/guodongq/quickstart/config"
	"github.com/guodongq/quickstart/pkg/provider/app"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"probes",
	fx.Provide(
		func(appProvider *app.App, cfg *config.Config) *Probes {
			probesCfg := cfg.Probes
			return New(
				appProvider,
				WithProbesOptionsLiveEndpoint(probesCfg.LivenessEndpoint),
				WithProbesOptionsReadEndpoint(probesCfg.ReadinessEndpoint),
				WithProbesOptionsPort(probesCfg.Port),
				WithProbesOptionsEnable(probesCfg.Enabled),
			)
		},
	),
	fx.Invoke(func(probesProvider *Probes, lc fx.Lifecycle) error {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go probesProvider.Run()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return probesProvider.Close()
			},
		})
		return probesProvider.Init()
	}),
)
