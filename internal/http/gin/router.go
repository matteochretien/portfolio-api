package gin

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/gin-gonic/gin"
)

type GinHttpRouter struct {
	httpServer *GinHttpServer
	ginGroup   gin.IRouter
}

func NewGinHttpRouter(httpServer *GinHttpServer, ginGroup gin.IRouter) *GinHttpRouter {
	return &GinHttpRouter{
		httpServer: httpServer,
		ginGroup:   ginGroup,
	}
}

func (r *GinHttpRouter) GET(path string, handlers ...http.HttpHandler) {
	r.ginGroup.GET(path, r.transformHandlers(handlers)...)
}

func (r *GinHttpRouter) POST(path string, handlers ...http.HttpHandler) {
	r.ginGroup.POST(path, r.transformHandlers(handlers)...)
}

func (r *GinHttpRouter) PUT(path string, handlers ...http.HttpHandler) {
	r.ginGroup.PUT(path, r.transformHandlers(handlers)...)
}

func (r *GinHttpRouter) DELETE(path string, handlers ...http.HttpHandler) {
	r.ginGroup.DELETE(path, r.transformHandlers(handlers)...)
}

func (r *GinHttpRouter) Group(path string) http.HttpRouter {
	newGinGroup := r.ginGroup.Group(path)
	return NewGinHttpRouter(r.httpServer, newGinGroup)
}

func (r *GinHttpRouter) Use(handlers ...http.HttpHandler) {
	r.ginGroup.Use(r.transformHandlers(handlers)...)
}

func (r *GinHttpRouter) transformHandlers(handlers []http.HttpHandler) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler // create a new variable to avoid closure over the same variable
		ginHandlers[i] = func(c *gin.Context) {
			handler(NewGinContext(c))
		}
	}

	return ginHandlers
}
