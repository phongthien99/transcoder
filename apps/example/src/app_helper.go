package src

import (
	"fmt"
	"time"

	"github.com/go-swagno/swagno"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/sigmaott/gest/package/extension/echofx/echo-swagger"
	"go.uber.org/zap"
)

func EnableSwagger(e *echo.Group, logger *zap.SugaredLogger, sw *swagno.Swagger) {
	sw.Host = ""
	e.GET("/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		f := func() []byte {
			return sw.MustToJson()
		}
		c.DynamicLoad = &f
	}))
}

func EnableLogRequest(e *echo.Group) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {

			before := v.StartTime
			after := time.Now()
			format := "\033[32m%s\033[0m \n"

			if c.Response().Status >= 300 {
				format = "\033[33m%s\033[0m \n"
			}
			if v.Status >= 400 {
				format = "\033[31m%s\033[0m"
			}
			message := fmt.Sprintf("[REQUEST] %v |%v|%v| %v| %v %v \n", after.Format(time.RFC3339), v.Status, after.Sub(before), c.Request().RemoteAddr, c.Request().Method, c.Request().URL)
			fmt.Printf(format, message)
			return nil
		},
	}))

}
