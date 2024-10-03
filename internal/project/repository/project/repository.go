package project

import (
	"context"
	"errors"
	"github.com/Niromash/niromash-api/internal/project/entity"
)

var ProjectNotFoundErr = errors.New("project not found")

type ProjectRepository interface {
	Add(ctx context.Context, project *entity.Project) error
	Get(ctx context.Context, projectId entity.ProjectId) (project *entity.Project, err error)
	List(ctx context.Context) (projects []*entity.Project, err error)
	Remove(ctx context.Context, projectId entity.ProjectId) error
	Update(ctx context.Context, project *entity.Project) error
}
