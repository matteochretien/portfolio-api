package list_projects

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/google/uuid"
	http2 "net/http"
	"time"
)

type ListProjectsHandler struct {
	listProjectsUseCase *ListProjectsUseCase
}

func NewListProjectsHandler(listProjectsUseCase *ListProjectsUseCase) *ListProjectsHandler {
	return &ListProjectsHandler{listProjectsUseCase}
}

func (h *ListProjectsHandler) Handle(c http.HttpContext) {
	projects, err := h.listProjectsUseCase.Execute(c)
	if err != nil {
		c.JSON(http2.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to list projects",
		})
		return
	}

	response := make([]listProjectResponse, len(projects))
	for i, project := range projects {
		response[i] = listProjectResponse{
			Id:              uuid.UUID(project.Id),
			Name:            project.Name,
			Description:     project.Description,
			PreviewImageUrl: project.PreviewImageUrl,
			LinkedLink:      project.LinkedLink,
			GithubLink:      project.GithubLink,
			Date:            project.Date,
			Client:          project.Client,
			ImagesLinks:     project.ImagesLinks,
			Categories:      make([]listProjectCategoryResponse, len(project.Categories)),
			Techs:           make([]listProjectTechResponse, len(project.TechStack)),
		}

		for j, cat := range project.Categories {
			response[i].Categories[j] = listProjectCategoryResponse{
				Id:   uuid.UUID(cat.Id),
				Name: cat.Name,
			}
		}

		for j, tech := range project.TechStack {
			response[i].Techs[j] = listProjectTechResponse{
				Id:   uuid.UUID(tech.Id),
				Name: tech.Name,
			}
		}
	}

	c.JSON(http2.StatusOK, response)
}

type listProjectResponse struct {
	Id              uuid.UUID                     `json:"id"`
	Name            string                        `json:"name"`
	Description     string                        `json:"description"`
	PreviewImageUrl string                        `json:"previewImageUrl"`
	LinkedLink      string                        `json:"linkedLink"`
	GithubLink      string                        `json:"githubLink"`
	Date            time.Time                     `json:"date"`
	Client          string                        `json:"client"`
	ImagesLinks     []string                      `json:"imagesLinks"`
	Categories      []listProjectCategoryResponse `json:"categories"`
	Techs           []listProjectTechResponse     `json:"techs"`
}

type listProjectCategoryResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type listProjectTechResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
