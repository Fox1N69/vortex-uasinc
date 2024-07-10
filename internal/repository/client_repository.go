package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"test-task/internal/models"
	"time"
)

type ClientRepository interface {
	Create(client *models.Client) (int64, error)
	ClientByID(id int64) (*models.Client, error)
	Update(id int64, updateParams map[string]interface{}) error
	Delete(id int64) error
	Clients() ([]models.Client, error)	
}

type clientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (cr *clientRepository) Create(client *models.Client) (int64, error) {
	const op = "repository.client.Create"

	query := `
	INSERT INTO clients (client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	var id int64
	err := cr.db.QueryRow(
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
		return 0, fmt.Errorf("%s %w Cloud not create client", op, err)
	}

	return id, nil
}

func (cr *clientRepository) ClientByID(id int64) (*models.Client, error) {
	const op = "repository.client.ClientID"

	query := `
		SELECT * from clients
		WHERE id = $1
	`

	var client models.Client
	err := cr.db.QueryRow(query, id).Scan(&client)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%s %w Cloud not get client", op, err)
	}

	return &client, nil
}

func (cr *clientRepository) Update(id int64, updateParams map[string]interface{}) error {
	const op = "repository.client.Update"

	if len(updateParams) == 0 {
		return fmt.Errorf("%s No updates provider", op)
	}

	setClauses := make([]string, 0, len(updateParams))
	args := make([]interface{}, 0, len(updateParams)+1)
	i := 1

	for column, value := range updateParams {
		setClauses = append(setClauses, fmt.Sprintf("%s = %d", column, i))
		args = append(args, value)
		i++
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE clients SET %s, updated_at = $%d WHERE id = $%d", setClause, i, i+1)
	args = append(args, time.Now(), id)

	_, err := cr.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s %w Cloud not update client", op, err)
	}

	return nil
}

func (cr *clientRepository) Delete(id int64) error {
	const op = "repository.client.Delete"

	query := `
	DELETE FROM clients 
	WHERE id = $1
	`

	_, err := cr.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s %w Cloud not delete client", op, err)
	}

	return nil
}

func (cr *clientRepository) Clients() ([]models.Client, error) {
	const op = "repository.client.Clients"

	query := `
		SELECT * FROM clients
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s %w Cloud not list clients", op, err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client)
		if err != nil {
			return nil, fmt.Errorf("%s %w Cloud not scan client", op, err)
		}
		clients = append(clients, client)
	}

	return clients, nil
}
