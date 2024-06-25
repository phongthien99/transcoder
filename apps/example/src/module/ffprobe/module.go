package ffprobe

import (
	"example/src/module/ffprobe/controller"
	"example/src/module/ffprobe/service.go"

	"github.com/sigmaott/gest/package/extension/echofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("health",
		fx.Provide(
			echofx.AsRoute(controller.NewProbeController),
			service.NewProbeService,
		),
	)
}
