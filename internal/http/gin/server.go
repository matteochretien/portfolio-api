package gin

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/gin-gonic/gin"
)

type GinHttpServer struct {
	http.HttpRouter
	engine           *gin.Engine
	listeningAddress string
}

func NewGinHttpServer(listeningAddress string) func() http.HttpServer {
	engine := gin.New()
	httpServer := &GinHttpServer{engine: engine, listeningAddress: listeningAddress}
	httpServer.HttpRouter = &GinHttpRouter{httpServer: httpServer, ginGroup: engine.Group("")}

	return func() http.HttpServer {
		return httpServer
	}
}

func (s *GinHttpServer) Start() error {
	return s.engine.Run(s.listeningAddress)
}
