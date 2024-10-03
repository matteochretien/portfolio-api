package category

import (
	"context"
	"github.com/Niromash/niromash-api/internal/postgres"
	"github.com/Niromash/niromash-api/internal/project/entity"
	"github.com/google/uuid"
)

type PostgresCategoryRepository struct {
	client *postgres.Client
}

type postgresCategory struct {
	Id   uuid.UUID
	Name string
}

func NewPostgresCategoryRepository(client *postgres.Client) CategoryRepository {
	return &PostgresCategoryRepository{client: client}
}

func (p *PostgresCategoryRepository) Add(ctx context.Context, category *entity.ProjectCategory) error {
	query := `INSERT INTO project_category (name) VALUES ($1) RETURNING id`
	return p.client.Client.QueryRow(ctx, query, category.Name).Scan(&category.Id)
}

func (p *PostgresCategoryRepository) AddIfNotExists(ctx context.Context, category *entity.ProjectCategory) error {
	query := `
		WITH ins AS (
			INSERT INTO project_category (name) VALUES ($1)
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		)
		SELECT id FROM ins
		UNION ALL
		SELECT id FROM project_category WHERE name = $1
		LIMIT 1;
	`
	err := p.client.Client.QueryRow(ctx, query, category.Name).Scan(&category.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresCategoryRepository) Get(ctx context.Context, categoryId entity.ProjectCategoryId) (*entity.ProjectCategory, error) {
	query := `SELECT id, name FROM project_category WHERE id = $1`

	var category postgresCategory
	err := p.client.Client.QueryRow(ctx, query, categoryId).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, err
	}

	return &entity.ProjectCategory{
		Id:   entity.ProjectCategoryId(category.Id),
		Name: category.Name,
	}, nil
}

func (p *PostgresCategoryRepository) List(ctx context.Context) (categories []*entity.ProjectCategory, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresCategoryRepository) Remove(ctx context.Context, categoryId entity.ProjectCategoryId) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresCategoryRepository) Update(ctx context.Context, category *entity.ProjectCategory) error {
	//TODO implement me
	panic("implement me")
}
