package middlewares

import (
	"github.com/Niromash/niromash-api/internal/http"
	http2 "net/http"
)

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (l *CorsMiddleware) Handle(c http.HttpContext) {
	c.SetHeader("Access-Control-Allow-Origin", "*")
	c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.SetHeader("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.SetHeader("Access-Control-Allow-Credentials", "true")

	if c.GetMethod() == http2.MethodOptions {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
