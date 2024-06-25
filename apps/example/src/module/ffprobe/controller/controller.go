package controller

import (
	"log"
	"net/http"

	"example/src/module/ffprobe/dto"
	"example/src/module/ffprobe/service.go"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/labstack/echo/v4"
	"github.com/sigmaott/gest/package/extension/echofx/common"
	"github.com/sigmaott/gest/package/extension/echofx/pipe"
)

// IMetricController defines the interface for the Metric controller.
type IProbeController interface {
	Probe()
	FastConvert()
}

// metricController implements the IMetricController interface.
type probeController struct {
	router  *echo.Group
	swagger *swagno.Swagger
	service service.IProbeService
}

// NewMetricController creates a new instance of metricController.
func NewProbeController(router *echo.Group, swagger *swagno.Swagger, service service.IProbeService) IProbeController {
	return &probeController{
		router:  router.Group("/prope"),
		swagger: swagger,
		service: service,
	}
}

// FastConvert implements IProbeController.
func (e *probeController) FastConvert() {
	e.swagger.AddEndpoint(endpoint.New(
		endpoint.POST,
		"prope/fast-convert",
		endpoint.WithTags("prope"),
	))

	body := common.Body[dto.FastConvertBody, dto.FastConvertBody]("body")
	e.router.POST("/fast-convert", func(c echo.Context) error {
		payload := body.Get(c)
		log.Print(payload)
		err := e.service.FastConvert(payload.Input, payload.Output, payload.Options)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, payload)
	}, pipe.UsePipes(body))
}

// Metric sets up the metrics endpoint and Swagger documentation.
func (e *probeController) Probe() {
	e.swagger.AddEndpoint(endpoint.New(
		endpoint.GET,
		"/prope",
		endpoint.WithTags("prope"),
	))
	// query := common.Query[dto.ProbeQuery, dto.ProbeQuery]("query")
	e.router.GET("", func(c echo.Context) error {
		q := c.QueryParam("input")
		res, err := e.service.GetMultimediaInfo(q)
		if err != nil {
			return err
		}
		c.Response().Header().Add("Content-Type", "application/json")
		return c.String(http.StatusOK, res)
	})
}
