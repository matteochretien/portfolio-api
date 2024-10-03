package gin

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) http.HttpContext {
	return &GinContext{c}
}

func (g *GinContext) SetHeader(key string, value string) {
	g.Context.Header(key, value)
}

func (g *GinContext) GetMethod() string {
	return g.Context.Request.Method
}
