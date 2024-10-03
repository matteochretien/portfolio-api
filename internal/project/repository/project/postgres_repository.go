package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/Niromash/niromash-api/OLD/api"
	"github.com/Niromash/niromash-api/internal/postgres"
	"github.com/Niromash/niromash-api/internal/project/entity"
	"github.com/Niromash/niromash-api/internal/project/repository/category"
	"github.com/Niromash/niromash-api/internal/project/repository/tech"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type PostgresProjectRepository struct {
	client             *postgres.Client
	categoryRepository category.CategoryRepository
	techRepository     tech.TechRepository
}

type postgresProject struct {
	Id              uuid.UUID
	Name            string
	Description     string
	PreviewImageUrl string
	LinkedLink      string
	GithubLink      string
	Date            time.Time
	Client          string
	ImagesLinks     []string
}

func NewPostgresProjectRepository(client *postgres.Client, categoryRepository category.CategoryRepository, techRepository tech.TechRepository) ProjectRepository {
	return &PostgresProjectRepository{client: client, categoryRepository: categoryRepository, techRepository: techRepository}
}

func (p *PostgresProjectRepository) Add(ctx context.Context, project *entity.Project) error {
	tx, err := p.client.Client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert project
	query := `INSERT INTO projects (name, description, preview_image_url, linked_link, github_link, date, client, images_links) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = tx.QueryRow(ctx, query, project.Name, project.Description, project.PreviewImageUrl, project.LinkedLink, project.GithubLink, project.Date, project.Client, project.ImagesLinks).Scan(&project.Id)
	if err != nil {
		return fmt.Errorf("failed to add project: %w", err)
	}

	// Insert categories
	for _, cat := range project.Categories {
		err = p.categoryRepository.AddIfNotExists(ctx, cat)
		if err != nil {
			return fmt.Errorf("failed to add category: %w", err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO project_category_projects (project_id, category_id) VALUES ($1, $2)`, project.Id, cat.Id)
		if err != nil {
			return fmt.Errorf("failed to add project category: %w", err)
		}
	}

	// Insert techs
	for _, t := range project.TechStack {
		err = p.techRepository.AddIfNotExists(ctx, t)
		if err != nil {
			return fmt.Errorf("failed to add tech: %w", err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO project_techs_projects (project_id, tech_id) VALUES ($1, $2)`, project.Id, t.Id)
		if err != nil {
			return fmt.Errorf("failed to add project tech: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (p *PostgresProjectRepository) Get(ctx context.Context, projectId entity.ProjectId) (*entity.Project, error) {
	tx, err := p.client.Client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var pp postgresProject
	err = tx.QueryRow(ctx, `SELECT id, name, description, preview_image_url, linked_link, github_link, date, client, images_links FROM projects WHERE id = $1`, projectId).Scan(&pp.Id, &pp.Name, &pp.Description, &pp.PreviewImageUrl, &pp.LinkedLink, &pp.GithubLink, &pp.Date, &pp.Client, &pp.ImagesLinks)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, api.ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	project := &entity.Project{
		Id:              entity.ProjectId(pp.Id),
		Name:            pp.Name,
		Description:     pp.Description,
		PreviewImageUrl: pp.PreviewImageUrl,
		LinkedLink:      pp.LinkedLink,
		GithubLink:      pp.GithubLink,
		Date:            pp.Date,
		Client:          pp.Client,
		ImagesLinks:     pp.ImagesLinks,
	}

	err = p.populateCategories(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to populate categories: %w", err)
	}

	err = p.populateTechs(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to populate techs: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return project, nil
}

func (p *PostgresProjectRepository) List(ctx context.Context) (projects []*entity.Project, err error) {
	tx, err := p.client.Client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT id, name, description, preview_image_url, linked_link, github_link, date, client, images_links FROM projects`)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pp postgresProject
		err = rows.Scan(&pp.Id, &pp.Name, &pp.Description, &pp.PreviewImageUrl, &pp.LinkedLink, &pp.GithubLink, &pp.Date, &pp.Client, &pp.ImagesLinks)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		project := &entity.Project{
			Id:              entity.ProjectId(pp.Id),
			Name:            pp.Name,
			Description:     pp.Description,
			PreviewImageUrl: pp.PreviewImageUrl,
			LinkedLink:      pp.LinkedLink,
			GithubLink:      pp.GithubLink,
			Date:            pp.Date,
			Client:          pp.Client,
			ImagesLinks:     pp.ImagesLinks,
		}

		err = p.populateCategories(ctx, project)
		if err != nil {
			return nil, fmt.Errorf("failed to populate categories: %w", err)
		}

		err = p.populateTechs(ctx, project)
		if err != nil {
			return nil, fmt.Errorf("failed to populate techs: %w", err)
		}

		projects = append(projects, project)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return projects, nil
}

func (p *PostgresProjectRepository) Remove(ctx context.Context, projectId entity.ProjectId) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresProjectRepository) Update(ctx context.Context, project *entity.Project) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresProjectRepository) populateCategories(ctx context.Context, project *entity.Project) error {
	rows, err := p.client.Client.Query(ctx, `SELECT category_id FROM project_category_projects WHERE project_id = $1`, project.Id)
	if err != nil {
		return fmt.Errorf("failed to get project categories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId uuid.UUID
		err = rows.Scan(&categoryId)
		if err != nil {
			return fmt.Errorf("failed to scan project category: %w", err)
		}

		cat, err := p.categoryRepository.Get(ctx, entity.ProjectCategoryId(categoryId))
		if err != nil {
			return fmt.Errorf("failed to get project category: %w", err)
		}

		project.Categories = append(project.Categories, cat)
	}

	return nil
}

func (p *PostgresProjectRepository) populateTechs(ctx context.Context, project *entity.Project) error {
	rows, err := p.client.Client.Query(ctx, `SELECT tech_id FROM project_techs_projects WHERE project_id = $1`, project.Id)
	if err != nil {
		return fmt.Errorf("failed to get project techs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var techId uuid.UUID
		err = rows.Scan(&techId)
		if err != nil {
			return fmt.Errorf("failed to scan project tech: %w", err)
		}

		t, err := p.techRepository.Get(ctx, entity.ProjectTechId(techId))
		if err != nil {
			return fmt.Errorf("failed to get project tech: %w", err)
		}

		project.TechStack = append(project.TechStack, t)
	}

	return nil
}
