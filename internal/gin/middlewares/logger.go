package middlewares

import (
	"github.com/Niromash/niromash-api/internal/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	logg logger.Logger
}

func NewLoggerMiddleware(logg logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logg: logg}
}

func (l *LoggerMiddleware) Execute() []gin.HandlerFunc {
	switch loggerType := l.logg.(type) {
	case *zap.SugaredLogger:
		desugar := loggerType.Desugar()

		return gin.HandlersChain{
			ginzap.GinzapWithConfig(desugar, &ginzap.Config{
				TimeFormat: "15:04:05 02/01/2006",
				SkipPaths:  []string{"/health"},
			}),
			ginzap.RecoveryWithZap(desugar, true),
		}

	default:
		return gin.HandlersChain{gin.Logger(), gin.Recovery()}
	}
}
