package logrus

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"logger",
	fx.Provide(New),
)
