package controller

import (
	"example/src/module/example-module/dto"
	"example/src/module/example-module/model"
	"log"
	"strconv"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/labstack/echo/v4"
	"github.com/sigmaott/gest/package/extension/echofx/common"
	"github.com/sigmaott/gest/package/extension/echofx/pipe"
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
		router:  router.Group("/examples"),
		swagger: swagger,
	}

}

func (e *exampleController) Create() {

	e.swagger.AddEndpoint(endpoint.New(
		endpoint.POST,
		"/examples",
		endpoint.WithTags("example"),
		endpoint.WithBody(dto.ExampleDto{}),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(model.Example{}, "OK", "200")}),
		endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		endpoint.WithConsume([]mime.MIME{mime.JSON}),
	))
	idParam := common.Param[int, int]("id")
	body := common.Body[dto.ExampleDto, dto.ExampleDto]("body")
	e.router.POST("", func(c echo.Context) error {

		log.Print(body.Get(c))
		// appIdHeader := common.Header[string, string]("x-app-id")
		// query := common.Query[dto.QueryExample, dto.QueryExample]("query")
		return c.JSON(201, model.Example{})
	}, pipe.UsePipes(idParam, body))

}

// FindOne implements IExampleController.
func (e *exampleController) FindOne() {

	e.swagger.AddEndpoint(endpoint.New(
		endpoint.GET,
		"/examples/{id}",
		endpoint.WithTags("example"),
		endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(model.Example{}, "OK", "200")}),
		endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		endpoint.WithConsume([]mime.MIME{mime.JSON}),
	))
	idParam := common.Param[string, int]("id")
	e.router.GET("/:id", func(c echo.Context) error {
		log.Print(idParam.Get(c))

		// body := common.Body[dto.CreateExample, dto.CreateExample]("body")
		// appIdHeader := common.Header[string, string]("x-app-id")
		// query := common.Query[dto.QueryExample, dto.QueryExample]("query")
		return c.JSON(201, model.Example{})
	}, pipe.UsePipes(idParam.Transform(func(value any, key string) (r any, e error) {
		return strconv.Atoi(value.(string))
	})))
}

// Paginate implements IExampleController.
func (e *exampleController) Paginate() {
	// panic("unimplemented")
}

// UpdateOne implements IExampleController.
func (e *exampleController) UpdateOne() {
	// panic("unimplemented")
}

// DeleteOne implements IExampleController.
func (e *exampleController) DeleteOne() {

	e.swagger.AddEndpoint(endpoint.New(
		endpoint.DELETE,
		"/examples/{id}",
		endpoint.WithTags("example"),
		endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
		endpoint.WithSuccessfulReturns([]response.Response{response.New(model.Example{}, "OK", "200")}),
		endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
		endpoint.WithConsume([]mime.MIME{mime.JSON}),
	))

	idParam := common.Param[string, int]("id")

	e.router.DELETE("/:id", func(c echo.Context) error {
		log.Print(idParam.Get(c))

		// body := common.Body[dto.CreateExample, dto.CreateExample]("body")
		// appIdHeader := common.Header[string, string]("x-app-id")
		// query := common.Query[dto.QueryExample, dto.QueryExample]("query")
		return c.JSON(201, model.Example{})
	}, pipe.UsePipes(idParam.Transform(func(value any, key string) (r any, e error) {
		return strconv.Atoi(value.(string))
	})))

}
