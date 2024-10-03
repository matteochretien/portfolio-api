package std

import (
	"github.com/Niromash/niromash-api/internal/http"
	http2 "net/http"
)

type Group struct {
	basePath   string
	httpRouter http.HttpRouter
	handlers   []*GroupHandler
}

type GroupHandler struct {
	Method  string
	Path    string
	Handler http.HttpHandler
}

type StandardHttpHandler struct {
	Method  string
	Path    string
	Handler http.HttpHandler
}

func (g *Group) GET(path string, handlers ...http.HttpHandler) {
	g.handle(http2.MethodGet, path, handlers)
}

func (g *Group) POST(path string, handlers ...http.HttpHandler) {
	g.handle(http2.MethodPost, path, handlers)
}

func (g *Group) PUT(path string, handlers ...http.HttpHandler) {
	g.handle(http2.MethodPut, path, handlers)
}

func (g *Group) DELETE(path string, handlers ...http.HttpHandler) {
	g.handle(http2.MethodDelete, path, handlers)
}

func (g *Group) Group(path string, handlers ...http.HttpHandler) *Group {
	newGroup := &Group{
		basePath:   g.basePath + path,
		httpRouter: g.httpRouter,
		handlers:   make([]*GroupHandler, len(g.handlers)),
	}
	copy(newGroup.handlers, g.handlers)
	newGroup.Use(handlers...)
	return newGroup
}

func (g *Group) Use(handlers ...http.HttpHandler) {
	for _, handler := range handlers {
		g.handlers = append(g.handlers, &GroupHandler{Handler: handler})
	}
}

func (g *Group) handle(method, path string, handlers []http.HttpHandler) {
	g.httpRouter.Handle(method, g.basePath+path, handlers...)
}

func (g *Group) combineHandlers(handlers ...http.HttpHandler) []http.HttpHandler {
	finalSize := len(g.handlers) + len(handlers)
	mergedHandlers := make([]http.HttpHandler, finalSize)
	for i, handler := range g.handlers {
		mergedHandlers[i] = handler.Handler
	}
	for i, handler := range handlers {
		mergedHandlers[i+len(g.handlers)] = handler
	}
	return mergedHandlers
}
