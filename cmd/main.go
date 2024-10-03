package main

import (
	"github.com/Niromash/niromash-api/internal/dig"
	"github.com/Niromash/niromash-api/internal/http"
	middlewares2 "github.com/Niromash/niromash-api/internal/http/middlewares"
	"github.com/Niromash/niromash-api/internal/http/routes"
	"github.com/Niromash/niromash-api/internal/http/std"
	"github.com/Niromash/niromash-api/internal/logger"
	"github.com/Niromash/niromash-api/internal/logger/zap"
	"github.com/Niromash/niromash-api/internal/postgres"
	"github.com/Niromash/niromash-api/internal/project/repository/category"
	"github.com/Niromash/niromash-api/internal/project/repository/project"
	"github.com/Niromash/niromash-api/internal/project/repository/tech"
	"github.com/Niromash/niromash-api/internal/project/usecase/get_project"
	"github.com/Niromash/niromash-api/internal/project/usecase/list_projects"
	"github.com/Niromash/niromash-api/internal/user/gin/middlewares"
	user_repository "github.com/Niromash/niromash-api/internal/user/repository"
	"github.com/Niromash/niromash-api/internal/user/usecase"
	"github.com/Niromash/niromash-api/internal/user/usecase/crypt"
	"github.com/Niromash/niromash-api/internal/user/usecase/get_account"
	"github.com/Niromash/niromash-api/internal/user/usecase/login"
	"github.com/Niromash/niromash-api/internal/user/usecase/register"
	"github.com/Niromash/niromash-api/internal/user/usecase/token"
	uberdig "go.uber.org/dig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var di = uberdig.New()

func main() {
	provideAll(di)

	var logg logger.Logger
	if err := di.Invoke(func(logger logger.Logger) {
		logg = logger
	}); err != nil {
		log.Fatalf("Unable to invoke : %s", err.Error())
	}

	invokeFatal(logg, di, func(httpServer http.HttpServer, corsMiddleware *middlewares2.CorsMiddleware, rtes *routes.Routes) {
		httpServer.Use(corsMiddleware.Handle)
		routes.RegisterRoutesForHttpServer(rtes, httpServer)
	})

	logg.Info("Starting API...")

	go invokeFatal(logg, di, runHttpServer(logg))

	logg.Info("API started!")

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, os.Interrupt, syscall.SIGTERM)

	<-signCh
	logg.Info("Shutting down GoRoutine...")
}

func invokeFatal(logg logger.Logger, di *uberdig.Container, f any) {
	if err := di.Invoke(f); err != nil {
		logg.Fatalf("Unable to invoke : %s", err.Error())
	}
}

func runHttpServer(logg logger.Logger) func(http.HttpServer) {
	return func(server http.HttpServer) {
		err := server.Start()
		if err != nil {
			logg.Fatal("Unable to run the router : %s", err.Error())
		}
	}
}

func provideAll(di *uberdig.Container) {
	providers := []dig.Provider{
		// Logger
		dig.NewProvider(zap.NewZapLogger),
		// Postgres
		dig.NewProvider(postgres.NewClient),
		// HttpServer
		dig.NewProvider(std.NewStandardHttpServer(os.Getenv("HTTP_SERVER_ADDRESS"))),
		//dig.NewProvider(gin2.NewGinHttpServer(os.Getenv("HTTP_SERVER_ADDRESS"))),
		dig.NewProvider(middlewares2.NewLoggerMiddleware),
		dig.NewProvider(middlewares2.NewCorsMiddleware),
		dig.NewProvider(http.NewHealthHandler),
		dig.NewProvider(routes.NewRoutes),
		// Gin
		//dig.NewProvider(gin.NewGinRouter(os.Getenv("HTTP_SERVER_PORT"))),
		//dig.NewProvider(gin_middlewares.NewApiKeyMiddleware),
		//dig.NewProvider(gin_middlewares.NewLoggerMiddleware),
		// User repositories
		dig.NewProvider(user_repository.NewPostgresUserRepository),
		// User usecases
		dig.NewProvider(usecase.NewGetUserByIdUseCase),
		dig.NewProvider(usecase.NewGetUserByEmailUseCase),
		dig.NewProvider(login.NewLoginUserUseCase),
		dig.NewProvider(usecase.NewExistsUserUseCase),
		dig.NewProvider(register.NewRegisterUserUseCase),
		dig.NewProvider(token.NewRefreshAccessTokenUseCase),
		// User handlers
		dig.NewProvider(get_account.NewGetUserHandler),
		dig.NewProvider(login.NewLoginHandler),
		dig.NewProvider(token.NewRefreshTokenHandler),
		dig.NewProvider(register.NewRegisterHandler),
		// User misc
		dig.NewProvider(token.NewJwtTokenStrategy(os.Getenv("JWT_SECRET"))),
		dig.NewProvider(crypt.NewBcryptStrategy),
		dig.NewProvider(middlewares.NewAuthMiddleware),
		// Projects repositories
		dig.NewProvider(project.NewPostgresProjectRepository),
		dig.NewProvider(tech.NewPostgresTechRepository),
		dig.NewProvider(category.NewPostgresCategoryRepository),
		// Projects usecases
		dig.NewProvider(get_project.NewGetProjectUseCase),
		dig.NewProvider(list_projects.NewListProjectsUseCase),
		// Projects handlers
		dig.NewProvider(get_project.NewGetProjectsHandler),
		dig.NewProvider(list_projects.NewListProjectsHandler),
	}
	for _, provider := range providers {
		if err := di.Provide(provider.Constructor, provider.ProvideOptions...); err != nil {
			log.Fatalf("Unable to provide %s : %s", provider, err.Error())
		}
	}
}
