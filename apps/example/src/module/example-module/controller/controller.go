package controller

import (
	"example/src/module/example-module/dto"
	"example/src/module/example-module/model"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/labstack/echo/v4"
)

type IExampleController interface {
	Create()
	FindOne()
	Paginate()
	UpdateOne()
	DeleteOne()
}

type exampleController struct {
	router  *echo.Group
	swagger *swagno.Swagger
}

func NewExampleController(router *echo.Group, swagger *swagno.Swagger) IExampleController {
	return &exampleController{
		router:  router.Group("/example"),
		swagger: swagger,
	}

}

func (e *exampleController) Create() {

	e.swagger.AddEndpoint(endpoint.New(
		endpoint.POST,
		"/m3u8-crawl",
		endpoint.WithTags("m3u8"),
		endpoint.WithBody(dto.ExampleDto{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(model.Example{}, "OK", "200")}),
		endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		endpoint.WithConsume([]mime.MIME{mime.JSON}),
	))
	e.router.POST("", func(c echo.Context) error {

		return c.JSON(201, model.Example{})
	})

}

// DeleteOne implements IExampleController.
func (e *exampleController) DeleteOne() {
	// panic("unimplemented")
}

// FindOne implements IExampleController.
func (e *exampleController) FindOne() {
	// panic("unimplemented")
}

// Paginate implements IExampleController.
func (e *exampleController) Paginate() {
	// panic("unimplemented")
}

// UpdateOne implements IExampleController.
func (e *exampleController) UpdateOne() {
	// panic("unimplemented")
}
