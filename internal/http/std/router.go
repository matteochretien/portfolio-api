package std

import (
	"fmt"
	"github.com/Niromash/niromash-api/internal/http"
)

type StandardHttpRouter struct {
	httpServer *StandardHttpServer
	group      *Group
}

func NewStandardHttpRouter(httpServer *StandardHttpServer, group *Group) *StandardHttpRouter {
	return &StandardHttpRouter{
		httpServer: httpServer,
		group:      group,
	}
}

func (r *StandardHttpRouter) GET(path string, handlers ...http.HttpHandler) {
	r.group.GET(path, handlers...)
}

func (r *StandardHttpRouter) POST(path string, handlers ...http.HttpHandler) {
	r.group.POST(path, handlers...)
}

func (r *StandardHttpRouter) PUT(path string, handlers ...http.HttpHandler) {
	r.group.PUT(path, handlers...)
}

func (r *StandardHttpRouter) DELETE(path string, handlers ...http.HttpHandler) {
	r.group.DELETE(path, handlers...)
}

func (r *StandardHttpRouter) Group(path string, handlers ...http.HttpHandler) http.HttpRouter {
	group := &Group{
		basePath:   path,
		httpRouter: r,
		handlers:   []*GroupHandler{},
	}
	handlers = group.combineHandlers(handlers...)
	group.Use(handlers...)

	return NewStandardHttpRouter(r.httpServer, group)
}

func (r *StandardHttpRouter) Use(handlers ...http.HttpHandler) {
	r.group.Use(handlers...)
}

func (r *StandardHttpRouter) Handle(method, path string, handlers ...http.HttpHandler) {
	if path == "" {
		return
	}
	fmt.Println("Handling", method, path, "with", len(handlers), "handlers")
	r.httpServer.handle(method, path, handlers...)
}
