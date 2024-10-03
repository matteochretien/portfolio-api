package get_project

import (
	"context"
	"github.com/Niromash/niromash-api/internal/project/entity"
	"github.com/Niromash/niromash-api/internal/project/repository/project"
)

type GetProjectUseCase struct {
	projectRepository project.ProjectRepository
}

func NewGetProjectUseCase(projectRepository project.ProjectRepository) *GetProjectUseCase {
	return &GetProjectUseCase{projectRepository: projectRepository}
}

func (u *GetProjectUseCase) Execute(ctx context.Context, projectId entity.ProjectId) (*entity.Project, error) {
	p, err := u.projectRepository.Get(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return p, nil
}
