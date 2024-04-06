package freestylelibre

import (
	"fmt"

	"github.com/TikhonP/maigo"
	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/TikhonP/medsenger-freestylelibre-bot/handler"
	"github.com/TikhonP/medsenger-freestylelibre-bot/util"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type handlers struct {
	root     *handler.RootHandler
	init     *handler.InitHandler
	status   *handler.StatusHandler
	remove   *handler.RemoveHandler
	settings *handler.SettingsHandler
}

func createHandlers(maigoClient *maigo.Client) *handlers {
	return &handlers{
		init: &handler.InitHandler{MaigoClient: maigoClient},
	}
}

func Serve(cfg *config.Server) {
	maigoClient := maigo.Init(cfg.MedsengerAgentKey)
	handlers := createHandlers(maigoClient)

	app := echo.New()
	app.Debug = cfg.Debug
	app.HideBanner = true
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: app.Logger.Output(),
	}))
	app.Use(middleware.Recover())
	//app.Use(middleware.BodyDump(func(context echo.Context, req []byte, res []byte) {
	//	fmt.Println()
	//	fmt.Println("Request:", string(req))
	//	fmt.Println("Response:", string(res))
	//}))
	if !cfg.Debug {
		app.Use(sentryecho.New(sentryecho.Options{}))
	}
	app.Validator = util.NewDefaultValidator()

	app.File("/styles.css", "public/styles.css")
	app.GET("/", handlers.root.Handle)
	app.POST("/init", handlers.init.Handle, util.ApiKeyJSON(cfg))
	app.POST("/status", handlers.status.Handle, util.ApiKeyJSON(cfg))
	app.POST("/remove", handlers.remove.Handle, util.ApiKeyJSON(cfg))

	app.GET("/settings", handlers.settings.Get, util.ApiKeyGetParam(cfg))
	app.POST("/settings", handlers.settings.Post, util.ApiKeyGetParam(cfg))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	err := app.Start(addr)
	if err != nil {
		panic(err)
	}
}
