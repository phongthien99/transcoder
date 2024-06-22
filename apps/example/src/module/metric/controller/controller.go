package controller

import (
	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// IMetricController defines the interface for the Metric controller.
type IMetricController interface {
	Metric()
}

// metricController implements the IMetricController interface.
type metricController struct {
	router   *echo.Group
	swagger  *swagno.Swagger
	registry *prometheus.Registry
}

// NewMetricController creates a new instance of metricController.
func NewMetricController(router *echo.Group, swagger *swagno.Swagger, registry *prometheus.Registry) IMetricController {
	return &metricController{
		router:   router.Group("/metrics"),
		swagger:  swagger,
		registry: registry,
	}
}

// Metric sets up the metrics endpoint and Swagger documentation.
func (e *metricController) Metric() {
	e.swagger.AddEndpoint(endpoint.New(
		endpoint.GET,
		"/metrics",
		endpoint.WithTags("metric"),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New("", "200", ""),
		}),
	))

	e.router.GET("", echo.WrapHandler(promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{
		Registry: e.registry,
	})))
}
