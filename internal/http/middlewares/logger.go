package middlewares

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/Niromash/niromash-api/internal/logger"
)

type LoggerMiddleware struct {
	logg logger.Logger
}

func NewLoggerMiddleware(logg logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logg: logg}
}

func (l *LoggerMiddleware) Handle(c http.HttpContext) {
	l.logg.Info("Request received")
	c.Next()
	l.logg.Info("Request processed")
}
