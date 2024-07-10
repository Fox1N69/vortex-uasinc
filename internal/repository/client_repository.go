package repository

import (
	"database/sql"
	"fmt"
	"test-task/internal/models"
)

type ClientRepository interface {
}

type clientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(client *models.Client) (int64, error) {
	const op = "repository.client.Create"

	query := `
	INSERT INTO clients (client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	var id int64
	err := r.db.QueryRow(
		query,
		client.ClientName,
		client.Version,
		client.Image,
		client.CPU,
		client.Memory,
		client.Priority,
		client.NeedRestart,
		client.SpawnedAt,
		client.CreatedAt,
		client.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Cloud not create client %s %w", op, err)
	}

	return id, nil
}
