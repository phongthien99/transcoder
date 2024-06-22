package metric

import (
	"example/src/module/metric/controller"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sigmaott/gest/package/extension/echofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("metric",
		fx.Provide(
			echofx.AsRoute(controller.NewMetricController),
			func() *prometheus.Registry {
				return prometheus.NewRegistry()
			},
		),
	)
}
