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
	logger zerolog.Logger
}

func NewServer(host string, port string, createCertification create.Certification, updateAnchor update.CertificationAnchor, maxMemoryPerRequest int64, logger zerolog.Logger) *Server {
	router := gin.Default()
	router.MaxMultipartMemory = maxMemoryPerRequest << 20

	v1 := router.Group("/v1/")
	v1.POST("certification", handler.PostCreateCertification(createCertification))

	v1.POST("webhook", handler.PostReceiveConfirmation(updateAnchor)) //TODO validacion webhook
	return &Server{host: host, port: port, engine: router, logger: logger}
}

func (s *Server) Start() error {
	if err := s.engine.Run(fmt.Sprintf("%s:%s", s.host, s.port)); err != nil {
		return err
	}
	return nil
}
