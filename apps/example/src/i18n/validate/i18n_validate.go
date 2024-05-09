package module

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	viTranslations "github.com/go-playground/validator/v10/translations/vi"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type I18nValidate struct {
	Ut        *ut.UniversalTranslator
	Validator *validator.Validate
}

func (i *I18nValidate) Validate(v any) error {
	if err := i.Validator.Struct(v); err != nil {
		return err
	}
	return nil
}

func NewI18nValidate(validate *validator.Validate, Ut *ut.UniversalTranslator, validateMessages map[string]map[string]string) *I18nValidate {

	for key, value := range validateMessages {
		trans, found := Ut.GetTranslator(key)
		if !found {
			continue
		}
		if key == "en" {
			enTranslations.RegisterDefaultTranslations(validate, trans)
		}
		if key == "vi" {
			viTranslations.RegisterDefaultTranslations(validate, trans)
		}
		for s, message := range value {
			validate.RegisterTranslation(s, trans, func(ut ut.Translator) error {
				return ut.Add(s, message, true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(s, fe.Field())
				return t
			})
		}

	}

	return &I18nValidate{
		Ut:        Ut,
		Validator: validate,
	}
}

func (i *I18nValidate) ValidateExceptionFilter() func(err error, c echo.Context) error {

	return func(err error, c echo.Context) error {
		trans, found := i.Ut.GetTranslator(i.GetAcceptLanguage(c))
		if !found {
			trans = i.Ut.GetFallback()
		}

		if he, ok := err.(validator.ValidationErrors); ok {
			errorBadRequest := map[string]any{
				"statusCode": http.StatusBadRequest,
				"message":    "Bad Request",
				"path":       c.Request().URL.Path,
				"timestamp":  time.Now().UnixMilli(),
				"errors": lo.MapToSlice(he.Translate(trans), func(key string, value string) string {

					return fmt.Sprintf("%s: %s", key[strings.Index(key, ".")+1:], value)
				}),
			}
			return c.JSON(http.StatusBadRequest, errorBadRequest)

		}
		return err

	}

}

func (i *I18nValidate) ValidatePipe() func(value any, key string) (r any, e error) {

	return func(value any, key string) (r any, e error) {
		err := i.Validate(value)
		return value, err
	}

}
func (i *I18nValidate) GetAcceptLanguage(c echo.Context) string {
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
	return i.Ut.GetFallback().Locale()
}
