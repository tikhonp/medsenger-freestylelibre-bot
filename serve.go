// Package freestylelibre is a web server for the Freestyle Libre integration bot.
package freestylelibre

import (
	"fmt"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tikhonp/maigo"
	"github.com/tikhonp/medsenger-freestylelibre-bot/handler"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

type Server struct {
	cfg      *util.Server
	client   *maigo.Client
	root     handler.RootHandler
	init     handler.InitHandler
	status   handler.StatusHandler
	remove   handler.RemoveHandler
	settings handler.SettingsHandler
}

func NewServer(cfg *util.Server) *Server {
	maigoClient := maigo.Init(cfg.MedsengerAgentKey)
	return &Server{
		cfg:    cfg,
		client: maigoClient,
		init:   handler.InitHandler{MaigoClient: maigoClient},
	}
}

func (s *Server) Listen() {
	app := echo.New()

	app.Debug = s.cfg.Debug
	app.HideBanner = true
	app.Validator = util.NewDefaultValidator()

	if !app.Debug {
		app.Use(sentryecho.New(sentryecho.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         5 * time.Second,
		}))
	}
	app.Use(middleware.RequestLoggerWithConfig(
		util.GetRequestLoggerConfig(!app.Debug),
	))
	app.Use(middleware.Recover())

	app.File("/styles.css", "public/styles.css")
	app.GET("/", s.root.Handle)
	app.POST("/init", s.init.Handle, util.AgentTokenJSON(s.client, maigo.RequestRoleSystem))
	app.POST("/status", s.status.Handle, util.AgentTokenJSON(s.client, maigo.RequestRoleSystem))
	app.POST("/remove", s.remove.Handle, util.AgentTokenJSON(s.client, maigo.RequestRoleSystem))

	app.GET("/settings", s.settings.Get, util.AgentTokenGetParam(s.client))
	app.POST("/settings", s.settings.Post, util.AgentTokenGetParam(s.client))

	app.GET("/setup", s.settings.Get, util.AgentTokenGetParam(s.client))
	app.POST("/setup", s.settings.Post, util.AgentTokenGetParam(s.client))

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	app.Logger.Fatal(app.Start(addr))
}
