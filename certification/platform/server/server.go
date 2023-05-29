package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	host   string
	port   string
	engine *gin.Engine
}

func (s Server) Start() error {
	router := gin.Default()

	return router.Run(fmt.Sprintf("%s:%s", s.host, s.port))
}

func (s Server) SetHandlers(f func(*gin.Engine)) {
	f(s.engine)
}
