package bus

import "go.uber.org/fx"

var Module = fx.Module(
	"bus",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(EventBus)),
		),
	),
)
