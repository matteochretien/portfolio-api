package gin

import (
	gin_middlewares "github.com/Niromash/niromash-api/internal/gin/middlewares"
	"github.com/Niromash/niromash-api/internal/user/gin/middlewares"
	"github.com/Niromash/niromash-api/internal/user/usecase/get_account"
	"github.com/Niromash/niromash-api/internal/user/usecase/login"
	"github.com/Niromash/niromash-api/internal/user/usecase/register"
	"github.com/Niromash/niromash-api/internal/user/usecase/token"
	"github.com/gin-gonic/gin"
	uberdig "go.uber.org/dig"
)

type GinRouter interface {
	RegisterRoutes()
	AddCorsHeaders()
	Run() error
}

type DefaultGinRouter struct {
	Router *gin.Engine
	deps   GinRouterDeps
	port   string
}

type GinRouterDeps struct {
	uberdig.In
	HealthHandler *HealthHandler

	// Middleware
	AuthMiddleware   *middlewares.AuthMiddleware
	ApiKeyMiddleware *gin_middlewares.ApiKeyMiddleware
	LoggerMiddleware *gin_middlewares.LoggerMiddleware
	// User
	GetUserHandler      *get_account.GetUserHandler
	LoginHandler        *login.LoginHandler
	RegisterHandler     *register.RegisterHandler
	RefreshTokenHandler *token.RefreshTokenHandler
}

func NewGinRouter(port string) func(deps GinRouterDeps) GinRouter {
	return func(deps GinRouterDeps) GinRouter {
		router := gin.New()
		return &DefaultGinRouter{Router: router, deps: deps, port: port}
	}
}

func (d *DefaultGinRouter) RegisterRoutes() {
	d.Router.Use(d.deps.LoggerMiddleware.Execute()...)
	authMiddleware := d.deps.AuthMiddleware.Execute

	d.Router.GET("/health", d.deps.HealthHandler.Handle)

	userGroup := d.Router.Group("/auth")
	{
		userGroup.GET("/me", authMiddleware, d.deps.GetUserHandler.Handle)
		userGroup.POST("/login", d.deps.LoginHandler.Handle)
		userGroup.POST("/register", d.deps.RegisterHandler.Handle)
		userGroup.POST("/token/access/renew", d.deps.RefreshTokenHandler.Handle)
	}
}

func (d *DefaultGinRouter) AddCorsHeaders() {
	d.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

func (d *DefaultGinRouter) Run() error {
	return d.Router.Run(":" + d.port)
}
