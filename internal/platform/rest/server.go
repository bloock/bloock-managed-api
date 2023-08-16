package rest

import (
	"bloock-managed-api/internal/platform/rest/handler"
	"bloock-managed-api/internal/service"
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

func NewServer(
	host string,
	port string,
	keysService service.GetLocalKeysService,
	managedKey service.ManagedKeyCreateService,
	localKey service.LocalKeyCreateService,
	signature service.SignService,
	createCertification service.CertificateService,
	updateAnchor service.CertificateUpdateAnchorService,
	webhookSecretKey string,
	enforceTolerance bool,
	logger zerolog.Logger,
	debug bool,
) (*Server, error) {
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
	v1.POST("sign", handler.PostSignData(signature))
	v1.POST("key/local", handler.PostCreateLocalKey(localKey))
	v1.POST("key/managed", handler.PostCreateManagedKey(managedKey))
	v1.GET("key/local", handler.GetLocalKeys(keysService))

	return &Server{host: host, port: port, engine: router, debug: debug, logger: logger}, nil
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
