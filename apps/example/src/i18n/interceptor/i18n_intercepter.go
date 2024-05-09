package interceptor

import (
	i18nError "example/src/i18n/error"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
)

type IInterceptor interface {
	Interceptor() func(c echo.Context, next echo.HandlerFunc) error
}

type II18nService interface {
	T(lang string, key string, params ...string) (string, error)

	C(lang string, key string, num float64, digits uint64, param string) (string, error)

	O(lang string, key string, num float64, digits uint64, param string) (string, error)

	R(lang string, key string, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error)
}

type I18nInterceptor struct {
	i18nService    II18nService
	fallbackLocale string
}

func (i *I18nInterceptor) Interceptor() func(c echo.Context, next echo.HandlerFunc) error {
	return func(c echo.Context, next echo.HandlerFunc) error {
		err := next(c)

		if err == nil {
			return nil

		}
		message, e := i.i18nService.T(i.GetAcceptLanguage(c), err.Error())
		log.Print(message, e)
		if e != nil {
			return i18nError.NewI18nError(err, "")

		}

		return i18nError.NewI18nError(err, message)

	}
}

func (i *I18nInterceptor) GetAcceptLanguage(c echo.Context) string {
	language := c.Request().Header.Get("Accept-Language")
	if language == "" {
		return ""
	}
	languageSplit := strings.Split(language, "-")
	if len(languageSplit) == 2 {
		return languageSplit[0]
	}
	if len(languageSplit) == 1 {
		return languageSplit[0]
	}
	return i.fallbackLocale
}

func NewI18nInterceter(i18nService II18nService, fallbackLocale string) IInterceptor {

	return &I18nInterceptor{
		i18nService:    i18nService,
		fallbackLocale: fallbackLocale,
	}

}
