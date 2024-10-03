package routes

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/Niromash/niromash-api/internal/project/usecase/get_project"
	"github.com/Niromash/niromash-api/internal/project/usecase/list_projects"
	uberdig "go.uber.org/dig"
)

type routesDeps struct {
	uberdig.In
	// Health
	HealthHandler *http.HealthHandler

	// Projects
	GetProjectHandler   *get_project.GetProjectHandler
	ListProjectsHandler *list_projects.ListProjectsHandler
}

type Routes struct {
	*routesDeps
}

func NewRoutes(deps routesDeps) *Routes {
	return &Routes{routesDeps: &deps}
}

func RegisterRoutesForHttpServer(routes *Routes, httpServer http.HttpServer) {
	httpServer.GET("/health", routes.HealthHandler.Handle)

	//projectsGroup := httpServer.Group("/projects")
	//projectsGroup.GET("/:id", routes.GetProjectHandler.Handle)
	//projectsGroup.GET("/", routes.ListProjectsHandler.Handle)
}
