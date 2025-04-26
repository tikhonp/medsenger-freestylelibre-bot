package freestylelibre

import (
	"fmt"

	"github.com/TikhonP/maigo"
	"github.com/tikhonp/medsenger-freestylelibre-bot/config"
	"github.com/tikhonp/medsenger-freestylelibre-bot/handler"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	cfg      *config.Server
	root     handler.RootHandler
	init     handler.InitHandler
	status   handler.StatusHandler
	remove   handler.RemoveHandler
	settings handler.SettingsHandler
}

func NewServer(cfg *config.Server) *Server {
	maigoClient := maigo.Init(cfg.MedsengerAgentKey)
	return &Server{
		cfg:  cfg,
		init: handler.InitHandler{MaigoClient: maigoClient},
	}
}

func (s *Server) Listen() {
	app := echo.New()
	app.Debug = s.cfg.Debug
	app.HideBanner = true
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: app.Logger.Output(),
	}))
	app.Use(middleware.Recover())
	if !s.cfg.Debug {
		app.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
		app.Logger.Printf("Sentry initialized")
	}
	app.Validator = util.NewDefaultValidator()

	app.File("/styles.css", "public/styles.css")
	app.GET("/", s.root.Handle)
	app.POST("/init", s.init.Handle, util.ApiKeyJSON(s.cfg))
	app.POST("/status", s.status.Handle, util.ApiKeyJSON(s.cfg))
	app.POST("/remove", s.remove.Handle, util.ApiKeyJSON(s.cfg))

	app.GET("/settings", s.settings.Get, util.ApiKeyGetParam(s.cfg))
	app.POST("/settings", s.settings.Post, util.ApiKeyGetParam(s.cfg))

	app.GET("/setup", s.settings.Get, util.ApiKeyGetParam(s.cfg))
	app.POST("/setup", s.settings.Post, util.ApiKeyGetParam(s.cfg))

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	app.Logger.Fatal(app.Start(addr))
}
