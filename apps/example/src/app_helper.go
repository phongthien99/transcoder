package src

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	i18nInterceptor "example/src/i18n/interceptor"
	i18nValidate "example/src/i18n/validate"
	validateMessage "example/src/locales/validate-message"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-swagno/swagno"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/sigmaott/gest/package/extension/echofx/echo-swagger"
	exceptionFilter "github.com/sigmaott/gest/package/extension/echofx/exception-filter"
	"github.com/sigmaott/gest/package/extension/echofx/interceptor"
	"github.com/sigmaott/gest/package/extension/i18nfx"
	"go.mongodb.org/mongo-driver/mongo"
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

func NewValidate(universalTranslator *ut.UniversalTranslator) *i18nValidate.I18nValidate {
	return i18nValidate.NewI18nValidate(validator.New(), universalTranslator, validateMessage.ValidateMessage)
}

func SetEchoInterceptor(e *echo.Echo, universalTranslator *ut.UniversalTranslator, i18nValidate *i18nValidate.I18nValidate, i18nService i18nfx.II18nService) {

	i18nInterceptorInstance := i18nInterceptor.NewI18nInterceter(i18nService, universalTranslator.GetFallback().Locale())
	e.Use(
		interceptor.UseInterceptors(i18nInterceptorInstance.Interceptor()),
		exceptionFilter.UseFilters(
			i18nValidate.ValidateExceptionFilter(),
			AllExceptionFilter,
		),
	)

}

func MongoNotFoundExceptionFilter(err error, c echo.Context) error {
	if errors.Is(err, mongo.ErrNoDocuments) {

		return c.JSON(http.StatusNotFound, map[string]any{
			"statusCode": 404,
			"message":    "Record Not Found",
			"path":       c.Request().URL.Path,
			"timestamp":  time.Now().UnixMilli(),
		})
	}
	return err
}

func AllExceptionFilter(err error, c echo.Context) error {

	if err != nil {
		e := err
		if errors.Unwrap(err) != nil {
			e = errors.Unwrap(err)
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"statusCode": 500,
			"message":    e.Error(),
			"path":       c.Request().URL.Path,
			"timestamp":  time.Now().UnixMilli(),
		})
	}

	return nil
}
