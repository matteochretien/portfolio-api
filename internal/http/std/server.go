package std

import (
	"fmt"
	"github.com/Niromash/niromash-api/internal/http"
	http2 "net/http"
)

type StandardHttpServer struct {
	http.HttpRouter
	engine           *http2.ServeMux
	listeningAddress string
}

func NewStandardHttpServer(listeningAddress string) func() http.HttpServer {
	engine := http2.NewServeMux()
	httpServer := &StandardHttpServer{engine: engine, listeningAddress: listeningAddress}
	httpServer.HttpRouter = &StandardHttpRouter{httpServer: httpServer, group: &Group{
		basePath:   "",
		httpRouter: httpServer,
		handlers:   []*GroupHandler{},
	}}

	return func() http.HttpServer {
		return httpServer
	}
}

func (s *StandardHttpServer) Start() error {
	return http2.ListenAndServe(s.listeningAddress, s.engine)
}

func (s *StandardHttpServer) handle(method, path string, handlers ...http.HttpHandler) {
	s.engine.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http2.ResponseWriter, r *http2.Request) {
		// Get handlers from groups
		// Merge them with the handlers from the route
		// Execute them

		//s.HttpRouter.(*StandardHttpRouter).group.handlers

		for _, handler := range handlers {
			handler(NewStandardContext(r, w))
		}
	})
}
