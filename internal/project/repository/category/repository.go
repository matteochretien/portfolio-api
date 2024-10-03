package category

import (
	"context"
	"errors"
	"github.com/Niromash/niromash-api/internal/project/entity"
)

var CategoryNotFoundErr = errors.New("category not found")

type CategoryRepository interface {
	Add(ctx context.Context, category *entity.ProjectCategory) error
	AddIfNotExists(ctx context.Context, category *entity.ProjectCategory) error
	Get(ctx context.Context, categoryId entity.ProjectCategoryId) (category *entity.ProjectCategory, err error)
	List(ctx context.Context) (categories []*entity.ProjectCategory, err error)
	Remove(ctx context.Context, categoryId entity.ProjectCategoryId) error
	Update(ctx context.Context, category *entity.ProjectCategory) error
}
