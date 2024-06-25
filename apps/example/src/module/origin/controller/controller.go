package controller

import (
	"example/src/module/origin/service"
	"io"
	"log"
	"net/http"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/labstack/echo/v4"
)

// IMetricController defines the interface for the Metric controller.
type IProbeController interface {
	Upload()
	Download()
}

// metricController implements the IMetricController interface.
type originController struct {
	router  *echo.Group
	swagger *swagno.Swagger
	service service.IFileUploadService
}

// NewMetricController creates a new instance of metricController.
func NewOriginController(router *echo.Group, swagger *swagno.Swagger, service service.IFileUploadService) IProbeController {
	return &originController{
		router:  router.Group("/origin"),
		swagger: swagger,
		service: service,
	}
}

// FastConvert implements IProbeController.
func (e *originController) Upload() {
	e.swagger.AddEndpoint(endpoint.New(
		endpoint.POST,
		"/origin/*",
		endpoint.WithTags("prope"),
	))

	e.router.PUT("/:config/:filePath", func(c echo.Context) error {
		uploadPath := c.Param("filePath")
		log.Printf("Upload path: %s", uploadPath)
		configBase64 := c.Param("config")
		defer c.Request().Body.Close()
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		e.service.UploadFile(uploadPath, body, configBase64)
		if err != nil {
			return err
		}
		log.Printf("File uploaded successfully")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "File uploaded successfully",
		})
	})
}

// Metric sets up the metrics endpoint and Swagger documentation.
func (e *originController) Download() {
	e.swagger.AddEndpoint(endpoint.New(
		endpoint.GET,
		"/prope",
		endpoint.WithTags("prope"),
	))
	// query := common.Query[dto.ProbeQuery, dto.ProbeQuery]("query")
	e.router.GET("/:config/:filePath", func(c echo.Context) error {

		filePath := c.Param("filePath")
		configBase64 := c.Param("config")

		content, err := e.service.ReadFile(filePath, configBase64)
		log.Printf("err %+v", err)
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "application/octet-stream", content)
	})
}
