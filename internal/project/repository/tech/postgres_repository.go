package tech

import (
	"context"
	"github.com/Niromash/niromash-api/internal/postgres"
	"github.com/Niromash/niromash-api/internal/project/entity"
	"github.com/google/uuid"
)

type PostgresTechRepository struct {
	client *postgres.Client
}

type postgresTech struct {
	Id   uuid.UUID
	Name string
}

func NewPostgresTechRepository(client *postgres.Client) TechRepository {
	return &PostgresTechRepository{client: client}
}

func (p *PostgresTechRepository) Add(ctx context.Context, tech *entity.ProjectTech) error {
	query := `INSERT INTO project_techs (name) VALUES ($1) RETURNING id`
	return p.client.Client.QueryRow(ctx, query, tech.Name).Scan(&tech.Id)
}

func (p *PostgresTechRepository) AddIfNotExists(ctx context.Context, tech *entity.ProjectTech) error {
	query := `
		WITH ins AS (
			INSERT INTO project_techs (name) VALUES ($1)
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		)
		SELECT id FROM ins
		UNION ALL
		SELECT id FROM project_techs WHERE name = $1
		LIMIT 1;
	`
	err := p.client.Client.QueryRow(ctx, query, tech.Name).Scan(&tech.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresTechRepository) Get(ctx context.Context, techId entity.ProjectTechId) (*entity.ProjectTech, error) {
	query := `SELECT id, name FROM project_techs WHERE id = $1`

	var tech postgresTech
	err := p.client.Client.QueryRow(ctx, query, techId).Scan(&tech.Id, &tech.Name)
	if err != nil {
		return nil, err
	}

	return &entity.ProjectTech{
		Id:   entity.ProjectTechId(tech.Id),
		Name: tech.Name,
	}, nil
}

func (p *PostgresTechRepository) List(ctx context.Context) (techs []*entity.ProjectTech, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresTechRepository) Remove(ctx context.Context, techId entity.ProjectTechId) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresTechRepository) Update(ctx context.Context, tech *entity.ProjectTech) error {
	//TODO implement me
	panic("implement me")
}
