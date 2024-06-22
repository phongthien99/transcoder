package health

import (
	"example/src/module/health/controller"
	"example/src/module/health/service"

	"github.com/sigmaott/gest/package/extension/echofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("health",
		fx.Provide(
			echofx.AsRoute(controller.NewMetricController),
			service.NewHeathCheckService,
		),
	)
}
