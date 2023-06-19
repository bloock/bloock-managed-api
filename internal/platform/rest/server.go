package rest

import (
	"bloock-managed-api/internal/platform/rest/handler"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/update"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	host   string
	port   string
	engine *gin.Engine
	debug  bool
	logger zerolog.Logger
}

func NewServer(host string, port string, createCertification create.Certification, updateAnchor update.CertificationAnchor, webhookSecretKey string, enforceTolerance bool, logger zerolog.Logger, debug bool) (*Server, error) {
	router := gin.Default()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	v1 := router.Group("/v1/")
	v1.POST("certification", handler.PostCreateCertification(createCertification))

	v1.POST("webhook", handler.PostReceiveWebhook(updateAnchor, webhookSecretKey, enforceTolerance))
	return &Server{host: host, port: port, engine: router, debug: debug, logger: logger}, nil
}

func (s *Server) Start() error {
	if err := s.engine.Run(fmt.Sprintf("%s:%s", s.host, s.port)); err != nil {
		return err
	}
	return nil
}
