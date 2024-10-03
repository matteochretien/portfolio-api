package get_project

import (
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/Niromash/niromash-api/internal/project/entity"
	"github.com/google/uuid"
	http2 "net/http"
	"time"
)

type GetProjectHandler struct {
	getProjectUseCase *GetProjectUseCase
}

func NewGetProjectsHandler(getProjectUseCase *GetProjectUseCase) *GetProjectHandler {
	return &GetProjectHandler{getProjectUseCase}
}

func (h *GetProjectHandler) Handle(c http.HttpContext) {
	projectId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http2.StatusBadRequest, map[string]any{
			"error":   err.Error(),
			"message": "Invalid project id",
		})
		return
	}

	project, err := h.getProjectUseCase.Execute(c, entity.ProjectId(projectId))
	if err != nil {
		c.JSON(http2.StatusInternalServerError, map[string]any{
			"error":   err.Error(),
			"message": "Failed to get project",
		})
		return
	}

	response := getProjectResponse{
		Id:              uuid.UUID(project.Id),
		Name:            project.Name,
		Description:     project.Description,
		PreviewImageUrl: project.PreviewImageUrl,
		LinkedLink:      project.LinkedLink,
		GithubLink:      project.GithubLink,
		Date:            project.Date,
		Client:          project.Client,
		ImagesLinks:     project.ImagesLinks,
		Categories:      make([]getProjectCategoryResponse, len(project.Categories)),
		Techs:           make([]getProjectTechResponse, len(project.TechStack)),
	}

	for j, cat := range project.Categories {
		response.Categories[j] = getProjectCategoryResponse{
			Id:   uuid.UUID(cat.Id),
			Name: cat.Name,
		}
	}

	for j, tech := range project.TechStack {
		response.Techs[j] = getProjectTechResponse{
			Id:   uuid.UUID(tech.Id),
			Name: tech.Name,
		}
	}

	c.JSON(http2.StatusOK, response)
}

type getProjectResponse struct {
	Id              uuid.UUID                    `json:"id"`
	Name            string                       `json:"name"`
	Description     string                       `json:"description"`
	PreviewImageUrl string                       `json:"previewImageUrl"`
	LinkedLink      string                       `json:"linkedLink"`
	GithubLink      string                       `json:"githubLink"`
	Date            time.Time                    `json:"date"`
	Client          string                       `json:"client"`
	ImagesLinks     []string                     `json:"imagesLinks"`
	Categories      []getProjectCategoryResponse `json:"categories"`
	Techs           []getProjectTechResponse     `json:"techs"`
}

type getProjectCategoryResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type getProjectTechResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
