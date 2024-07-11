package repository_test

import (
	"context"
	"test-task/internal/models"
	"test-task/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// TestCreate tests the creation of a client and its algorithm status in the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct insertion
// of a new client and its algorithm status records, and checks if the generated client ID matches
// the expected value.
func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewClientRepository(db)

	client := &models.Client{
		ClientName:  "TestClient",
		Version:     1,
		Image:       "test/image",
		CPU:         "2",
		Memory:      "1024",
		Priority:    1,
		NeedRestart: true,
		SpawnedAt:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	algorithm := &models.AlgorithmStatus{
		VWAP: false,
		TWAP: false,
		HFT:  false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO clients").
		WithArgs(client.ClientName, client.Version, client.Image, client.CPU, client.Memory, client.Priority, client.NeedRestart, client.SpawnedAt, client.CreatedAt, client.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("INSERT INTO algorithm_status").
		WithArgs(1, algorithm.VWAP, algorithm.TWAP, algorithm.HFT).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	id, err := repo.Create(client, algorithm)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestClientByID tests fetching a client by its ID from the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct retrieval
// of a client record by its ID and compares the fetched client with the expected client.
func TestClientByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewClientRepository(db)

	expectedClient := &models.Client{
		ID:          1,
		ClientName:  "TestClient",
		Version:     1,
		Image:       "test/image",
		CPU:         "2",
		Memory:      "1024",
		Priority:    1,
		NeedRestart: true,
		SpawnedAt:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery("SELECT \\* from clients WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "client_name", "version", "image", "cpu", "memory", "priority", "need_restart", "spawned_at", "created_at", "updated_at"}).
				AddRow(expectedClient.ID, expectedClient.ClientName, expectedClient.Version, expectedClient.Image, expectedClient.CPU, expectedClient.Memory, expectedClient.Priority, expectedClient.NeedRestart, expectedClient.SpawnedAt, expectedClient.CreatedAt, expectedClient.UpdatedAt))

	client, err := repo.ClientByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedClient, client)
	mock.ExpectationsWereMet()
}

// TestUpdate tests updating client information in the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct execution
// of an update query to modify client details based on provided update parameters.
func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewClientRepository(db)

	updateParams := map[string]interface{}{
		"client_name": "UpdatedClient",
		"priority":    2,
	}

	mock.ExpectExec("UPDATE clients SET client_name = \\$1, priority = \\$2, updated_at = \\$3 WHERE id = \\$4").
		WithArgs("UpdatedClient", 2, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(1, updateParams)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

// TestDelete tests deleting a client from the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct execution
// of a delete query to remove a client record by its ID.
func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewClientRepository(db)

	mock.ExpectExec("DELETE FROM clients WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

// TestClients tests fetching a list of clients from the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct retrieval
// of multiple client records and compares them with the expected client data.
func TestClients(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	_, err = rdb.Ping(context.Background()).Result()
	assert.NoError(t, err, "failed to ping Redis")

	repo := repository.NewClientRepository(db)

	ctx := context.Background()

	expectedClients := []models.Client{
		{
			ID:          1,
			ClientName:  "Client1",
			Version:     1,
			Image:       "image1",
			CPU:         "2",
			Memory:      "1024",
			Priority:    1,
			NeedRestart: false,
			SpawnedAt:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			ClientName:  "Client2",
			Version:     2,
			Image:       "image2",
			CPU:         "4",
			Memory:      "2048",
			Priority:    2,
			NeedRestart: true,
			SpawnedAt:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "client_name", "version", "image", "cpu", "memory", "priority", "need_restart", "spawned_at", "created_at", "updated_at"})
	for _, client := range expectedClients {
		rows.AddRow(client.ID, client.ClientName, client.Version, client.Image, client.CPU, client.Memory, client.Priority, client.NeedRestart, client.SpawnedAt, client.CreatedAt, client.UpdatedAt)
	}

	mock.ExpectQuery("SELECT id, client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at FROM clients").
		WillReturnRows(rows)

	clients, err := repo.Clients(ctx)
	assert.NoError(t, err, "error retrieving clients")

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, len(expectedClients), len(clients), "number of returned clients mismatch")
	for i := range expectedClients {
		assert.Equal(t, expectedClients[i], clients[i], "client mismatch")
	}
}

// TestAlgorithmStatuses tests fetching algorithm statuses from the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct retrieval
// of algorithm status records and compares them with the expected algorithm status data.
func TestAlgorithmStatuses(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewClientRepository(db)

	expectedStatuses := []models.AlgorithmStatus{
		{
			ID:       1,
			ClientID: 1,
			VWAP:     true,
			TWAP:     false,
			HFT:      true,
		},
		{
			ID:       2,
			ClientID: 2,
			VWAP:     false,
			TWAP:     true,
			HFT:      false,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "client_id", "vwap", "twap", "hft"})
	for _, status := range expectedStatuses {
		rows.AddRow(status.ID, status.ClientID, status.VWAP, status.TWAP, status.HFT)
	}

	mock.ExpectQuery("SELECT \\* from algorithm_status").
		WillReturnRows(rows)

	statuses, err := repo.AlgorithmStatuses()
	assert.NoError(t, err)
	assert.Equal(t, expectedStatuses, statuses)
	mock.ExpectationsWereMet()
}

// TestUpdateAlgorithmStatus tests updating algorithm status in the database.
//
// It mocks SQL database interactions using sqlmock. The test verifies the correct execution
// of an update query to modify algorithm status based on provided update parameters.
func TestUpdateAlgorithmStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewClientRepository(db)

	updateParams := map[string]interface{}{
		"vwap": true,
	}

	mock.ExpectExec("UPDATE algorithm_status SET vwap = \\$1 WHERE id = \\$2").
		WithArgs(true, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateAlgorithmStatus(1, updateParams)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}
