package list_projects

import (
	"context"
	"github.com/Niromash/niromash-api/internal/project/entity"
	"github.com/Niromash/niromash-api/internal/project/repository/project"
)

type ListProjectsUseCase struct {
	projectRepository project.ProjectRepository
}

func NewListProjectsUseCase(projectRepository project.ProjectRepository) *ListProjectsUseCase {
	return &ListProjectsUseCase{projectRepository: projectRepository}
}

func (u *ListProjectsUseCase) Execute(ctx context.Context) ([]*entity.Project, error) {
	projects, err := u.projectRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
