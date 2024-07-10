package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"test-task/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

type ClientRepository interface {
	Create(client *models.Client) (int64, error)
	ClientByID(id int64) (*models.Client, error)
	Update(id int64, updateParams map[string]interface{}) error
	Delete(id int64) error
	Clients(ctx context.Context) ([]models.Client, error)
	CreateAlgorithm(algorithm *models.AlgorithmStatus) (int64, error)
	AlgorithmStatuses() ([]models.AlgorithmStatus, error)
	AlgorithmByClientID(ctx context.Context, clientID int64) (*models.AlgorithmStatus, error)
	UpdateAlgorithmStatus(id int64, status map[string]interface{}) error
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
		return fmt.Errorf("%s No updates provided", op)
	}

	setClauses := make([]string, 0, len(updateParams))
	args := make([]interface{}, 0, len(updateParams)+2)
	i := 1

	for column, value := range updateParams {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
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

func (cr *clientRepository) Clients(ctx context.Context) ([]models.Client, error) {
	const op = "repository.client.Clients"

	query := `
		SELECT id, client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at
		FROM clients
	`

	rows, err := cr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s %w Cloud not list clients", op, err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client.ID, &client.ClientName, &client.Version, &client.Image, &client.CPU, &client.Memory, &client.Priority, &client.NeedRestart, &client.SpawnedAt, &client.CreatedAt, &client.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s %w Cloud not scan client", op, err)
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return clients, nil
}

func (cr *clientRepository) CreateAlgorithm(algorithm *models.AlgorithmStatus) (int64, error) {
	const op = "repository.client.CreateAlgorithm"

	tx, err := cr.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("%s: transaction rolled back due to error: %v", op, err)
		}
	}()

	query := `
	INSERT INTO algorithm_status (client_id, vwap, twap, hft)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var id int64
	err = tx.QueryRow(query, algorithm.ClientID, algorithm.VWAP, algorithm.TWAP, algorithm.HFT).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	logrus.Printf("%s: successfully created algorithm with ID %d", op, id)

	return id, nil
}

func (cr *clientRepository) AlgorithmStatuses() ([]models.AlgorithmStatus, error) {
	const op = "repository.client.AlgorithmStatuses"

	query := `
		SELECT * from algorithm_status
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s %w Cloud not list clients", op, err)
	}
	defer rows.Close()

	var statuses []models.AlgorithmStatus
	for rows.Next() {
		var status models.AlgorithmStatus
		err := rows.Scan(&status)
		if err != nil {
			return nil, fmt.Errorf("%s %w Clound not scan algorithm", op, err)
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (cr *clientRepository) UpdateAlgorithmStatus(id int64, status map[string]interface{}) error {
	const op = "repository.client.UpdateAlgorithmStatus"

	if len(status) == 0 {
		return fmt.Errorf("%s No updates provider", op)
	}

	setClauses := make([]string, 0, len(status))
	args := make([]interface{}, 0, len(status)+1)
	i := 1

	for column, value := range status {
		switch v := value.(type) {
		case bool:
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
			args = append(args, v)
		default:
			return fmt.Errorf("%s Unsupported type for column %s: %T", op, column, v)
		}
		i++
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE algorithm_status SET %s WHERE id = $%d", setClause, i)
	args = append(args, id)

	_, err := cr.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s %w Cloud not update algorithm", op, err)
	}

	return nil
}

func (cr *clientRepository) AlgorithmByClientID(ctx context.Context, clientID int64) (*models.AlgorithmStatus, error) {
	const op = "repository.client.AlgorithmByClientID"

	query := `
		SELECT id, client_id, vwap, twap, hft 
		FROM algorithm_status
		WHERE client_id = $1
	`

	var algorithm models.AlgorithmStatus
	err := cr.db.QueryRowContext(ctx, query, clientID).Scan(&algorithm.ID, &algorithm.ClientID, &algorithm.VWAP, &algorithm.TWAP, &algorithm.HFT)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &algorithm, nil
}
