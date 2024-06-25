package origin

import (
	"example/src/module/origin/controller"
	"example/src/module/origin/service"

	"github.com/sigmaott/gest/package/extension/echofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("health",
		fx.Provide(
			echofx.AsRoute(controller.NewOriginController),
			service.NewFileUploadService,
		),
	)
}
