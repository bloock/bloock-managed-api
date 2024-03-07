package rest

import (
	"fmt"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/handler"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/handler/process"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/handler/webhook"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/middleware"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	host   string
	port   string
	engine *gin.Engine
	logger zerolog.Logger
}

func NewServer(l zerolog.Logger) (*Server, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "gin").Logger()
	gin.DefaultWriter = l.With().Str("level", "info").Logger()
	gin.DefaultErrorWriter = l.With().Str("level", "error").Logger()

	router := gin.Default()
	if config.Configuration.Api.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	router.Use(middleware.ErrorMiddleware())
	router.Use(logger.SetLogger(
		logger.WithSkipPath([]string{"/health"}),
		logger.WithUTC(true),
		logger.WithLogger(func(c *gin.Context, _ zerolog.Logger) zerolog.Logger {
			return l
		}),
	))

	v1 := router.Group("/v1/")
	v1.GET("health", handler.Health())
	v1.POST("process", middleware.AuthMiddleware(), process.PostProcess(l))
	v1.POST("webhook", webhook.PostReceiveWebhook(l))
	v1.StaticFile("docs", "./internal/platform/rest/static/index.html")
	if config.Configuration.Api.DebugMode {
		v1.POST("debug", handler.Debug())
	}

	return &Server{host: config.Configuration.Api.Host, port: config.Configuration.Api.Port, engine: router, logger: l}, nil
}

func (s *Server) Start() error {
	if err := s.engine.Run(fmt.Sprintf("%s:%s", s.host, s.port)); err != nil {
		return err
	}
	return nil
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
