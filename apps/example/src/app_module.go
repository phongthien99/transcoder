package src

import (
	"example/src/config/types"
	example_http "example/src/module/example-module"
	"example/src/module/health"
	"example/src/module/metric"

	"example/src/locales/loader"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/vi"
	"github.com/go-swagno/swagno"
	"github.com/labstack/echo/v4"
	"github.com/sigmaott/gest/package/extension/configfx"
	"github.com/sigmaott/gest/package/extension/echofx"
	"github.com/sigmaott/gest/package/extension/i18nfx"
	"github.com/sigmaott/gest/package/extension/logfx"
	"github.com/sigmaott/gest/package/extension/metricfx"
	"github.com/sigmaott/gest/package/extension/swaggerfx"
	"go.uber.org/fx"
)

func NewApp(config *types.EnvironmentVariable) *fx.App {
	return fx.New(
		//config
		configfx.ForRoot(configfx.Option[*types.EnvironmentVariable]{
			Variable: config,
		}),
		// http
		echofx.ForRoot(config.Http.Port),
		fx.Provide(
			fx.Annotate(
				SetGlobalPrefix,
				fx.ParamTags(`name:"platformEcho"`),
			)),

		fx.Invoke(func(*echo.Echo) {}),
		// log
		logfx.ForRoot(config.Log.Level),
		// i18n
		i18nfx.ForRoot(i18nfx.I18nModuleParams{
			FallbackLanguage: "en",
			Loader:           loader.NewI18nMemoryLoader(),
			Translators: []locales.Translator{
				en.New(),
				vi.New(),
			},
		}),
		// i18n validate
		fx.Provide(fx.Annotate(
			NewValidate,
			fx.ParamTags(`name:"universalTranslator"`),
		),
		),

		// swagger
		swaggerfx.ForRoot(swagno.Config{
			Path: config.Http.BasePath,
		}),
		// metric
		metricfx.ForRoot(),
		// helper
		fx.Invoke(EnableSwagger),
		fx.Invoke(fx.Annotate(
			SetEchoInterceptor,
			fx.ParamTags(`name:"platformEcho"`, `name:"universalTranslator"`),
		)),
		// log http
		fx.Invoke(EnableLogRequest),

		// my module
		example_http.Module(),
		metric.Module(),
		health.Module(),
	)

}

func SetGlobalPrefix(e *echo.Echo, cfg configfx.IConfig[*types.EnvironmentVariable]) *echo.Group {

	// i18nInterceptorInstance := i18nInterceptor.NewI18nInterceter(i18nService, universalTranslator.GetFallback().Locale())
	// e.Use(
	// 	interceptor.UseInterceptors(i18nInterceptorInstance.Interceptor()),
	// 	exceptionFilter.UseFilters(
	// 		i18nValidate.ValidateExceptionFilter(),
	// 		sigmaCommon.AllExceptionFilter,
	// 	),
	// )
	return e.Group(cfg.Get().Http.BasePath)
}
