package tech

import (
	"context"
	"errors"
	"github.com/Niromash/niromash-api/internal/project/entity"
)

var TechNotFoundErr = errors.New("tech not found")

type TechRepository interface {
	Add(ctx context.Context, tech *entity.ProjectTech) error
	AddIfNotExists(ctx context.Context, tech *entity.ProjectTech) error
	Get(ctx context.Context, techId entity.ProjectTechId) (tech *entity.ProjectTech, err error)
	List(ctx context.Context) (techs []*entity.ProjectTech, err error)
	Remove(ctx context.Context, techId entity.ProjectTechId) error
	Update(ctx context.Context, tech *entity.ProjectTech) error
}
