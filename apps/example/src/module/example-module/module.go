package example_http

import (
	"example/src/module/example-module/controller"
	"example/src/module/example-module/repository"
	"example/src/module/example-module/service"

	"github.com/sigmaott/gest/package/extension/echofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("example",
		fx.Provide(
			echofx.AsRoute(controller.NewExampleController),
		),
		fx.Provide(
			service.NewExampleService,
		),
		fx.Provide(
			repository.NewExampleRepository,
		),
	)
}
