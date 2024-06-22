package controller

import (
	"example/src/module/health/service"
	"net/http"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/labstack/echo/v4"
)

// IMetricController defines the interface for the Metric controller.
type IHealthController interface {
	Health()
}

// metricController implements the IMetricController interface.
type metricController struct {
	router  *echo.Group
	swagger *swagno.Swagger
	service service.IHeathCheckService
}

// NewMetricController creates a new instance of metricController.
func NewMetricController(router *echo.Group, swagger *swagno.Swagger, service service.IHeathCheckService) IHealthController {
	return &metricController{
		router:  router.Group("/health"),
		swagger: swagger,
		service: service,
	}
}

// Metric sets up the metrics endpoint and Swagger documentation.
func (e *metricController) Health() {
	e.swagger.AddEndpoint(endpoint.New(
		endpoint.GET,
		"/health",
		endpoint.WithTags("health"),
		endpoint.WithSuccessfulReturns([]response.Response{
			response.New("", "200", ""),
		}),
	))

	e.router.GET("", func(c echo.Context) error {
		res, err := e.service.HeathCheckPayment()
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, res)
		}
		return c.JSON(http.StatusOK, res)
	})
}
