package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"test-task/internal/models"
	"time"
)

type ClientRepository interface {
	Create(client *models.Client, algorithm *models.AlgorithmStatus) (int64, error)
	ClientByID(id int64) (*models.Client, error)
	Update(id int64, updateParams map[string]interface{}) error
	Delete(id int64) error
	Clients(ctx context.Context) ([]models.Client, error)
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

// Create creates a new client record along with its associated algorithm status.
// It uses a transaction to ensure atomicity and returns the ID of the newly created client.
func (cr *clientRepository) Create(client *models.Client, algorithm *models.AlgorithmStatus) (int64, error) {
	const op = "repository.client.Create"

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

	queryClient := `
		INSERT INTO clients (client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	stmtClient, err := tx.Prepare(queryClient)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmtClient.Close()

	var clientID int64
	err = stmtClient.QueryRow(
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
	).Scan(&clientID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	queryAlgorithm := `
		INSERT INTO algorithm_status (client_id, vwap, twap, hft)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	stmtAlgorithm, err := tx.Prepare(queryAlgorithm)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmtAlgorithm.Close()

	algorithm.ClientID = clientID
	err = stmtAlgorithm.QueryRow(
		algorithm.ClientID,
		algorithm.VWAP,
		algorithm.TWAP,
		algorithm.HFT,
	).Scan(&algorithm.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return clientID, nil
}

// ClientByID retrieves a client by its ID from the database.
// It returns a pointer to the client object if found, or nil if not found.
func (cr *clientRepository) ClientByID(id int64) (*models.Client, error) {
	const op = "repository.client.ClientID"

	query := `
		SELECT id, client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at
		FROM clients
		WHERE id = $1
	`

	var client models.Client
	err := cr.db.QueryRow(query, id).Scan(
		&client.ID,
		&client.ClientName,
		&client.Version,
		&client.Image,
		&client.CPU,
		&client.Memory,
		&client.Priority,
		&client.NeedRestart,
		&client.SpawnedAt,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%s %w Could not get client", op, err)
	}

	return &client, nil
}

// Update updates a client record identified by the given ID with the provided update parameters.
// It accepts a map of update parameters where keys represent column names
// and values represent new values for those columns.
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
		return fmt.Errorf("%s %w Could not update client", op, err)
	}

	return nil
}

// Delete deletes a client record identified by the given ID from the database.
func (cr *clientRepository) Delete(id int64) error {
	const op = "repository.client.Delete"

	query := `
	DELETE FROM clients 
	WHERE id = $1
	`

	_, err := cr.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s %w Could not delete client", op, err)
	}

	return nil
}

// Clients retrieves all clients stored in the database.
// It returns a slice of client objects or an error if the operation fails.
func (cr *clientRepository) Clients(ctx context.Context) ([]models.Client, error) {
	const op = "repository.client.Clients"

	query := `
		SELECT id, client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at
		FROM clients
	`

	rows, err := cr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		err := rows.Scan(
			&client.ID,
			&client.ClientName,
			&client.Version,
			&client.Image,
			&client.CPU,
			&client.Memory,
			&client.Priority,
			&client.NeedRestart,
			&client.SpawnedAt,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return clients, nil
}

// AlgorithmStatuses retrieves all algorithm statuses stored in the database.
// It returns a slice of algorithm status objects or an error if the operation fails.
func (cr *clientRepository) AlgorithmStatuses() ([]models.AlgorithmStatus, error) {
	const op = "repository.client.AlgorithmStatuses"

	query := `
		SELECT id, client_id, vwap, twap, hft
		FROM algorithm_status
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s %w Could not list algorithm statuses", op, err)
	}
	defer rows.Close()

	var statuses []models.AlgorithmStatus
	for rows.Next() {
		var status models.AlgorithmStatus
		err := rows.Scan(
			&status.ID,
			&status.ClientID,
			&status.VWAP,
			&status.TWAP,
			&status.HFT,
		)
		if err != nil {
			return nil, fmt.Errorf("%s %w Could not scan algorithm status", op, err)
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// AlgorithmByClientID retrieves the algorithm status associated with a client ID.
// It returns a pointer to the algorithm status object if found, or nil if not found.
func (cr *clientRepository) UpdateAlgorithmStatus(id int64, status map[string]interface{}) error {
	const op = "repository.client.UpdateAlgorithmStatus"

	if len(status) == 0 {
		return fmt.Errorf("%s No updates provided", op)
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
		return fmt.Errorf("%s %w Could not update algorithm status", op, err)
	}

	return nil
}

// UpdateAlgorithmStatus updates the algorithm status identified by the given ID.
// It accepts a map of status updates where keys represent column names in the
// algorithm_status table and values represent new values for those columns.
func (cr *clientRepository) AlgorithmByClientID(ctx context.Context, clientID int64) (*models.AlgorithmStatus, error) {
	const op = "repository.client.AlgorithmByClientID"

	query := `
		SELECT id, client_id, vwap, twap, hft 
		FROM algorithm_status
		WHERE client_id = $1
	`

	var algorithm models.AlgorithmStatus
	err := cr.db.QueryRowContext(ctx, query, clientID).Scan(
		&algorithm.ID,
		&algorithm.ClientID,
		&algorithm.VWAP,
		&algorithm.TWAP,
		&algorithm.HFT,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &algorithm, nil
}
