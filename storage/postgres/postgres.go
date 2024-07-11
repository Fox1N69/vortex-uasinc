package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type PSQLClient struct {
	DB *sql.DB
}

func NewPSQLClient() *PSQLClient {
	return &PSQLClient{}
}

// Connect establishes a connection to the PostgreSQL database using the provided credentials.
//
// It sets up connection pooling with maximum open and idle connections,
// and sets the maximum lifetime of connections.
//
// Parameters:
// - user: PostgreSQL username.
// - password: Password for the PostgreSQL user.
// - host: PostgreSQL server host address.
// - port: PostgreSQL server port.
// - dbname: Name of the PostgreSQL database to connect to.
//
// Returns an error if the connection cannot be established or if there is an issue
// with setting up the connection pool or pinging the database.
func (s *PSQLClient) Connect(user, password, host, port, dbname string) error {
	const op = "storage.postgres.Connect()"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("%s: %w", op, err)
	}

	s.DB = db
	return nil
}

// Close closes the connection to the PostgreSQL database.
//
// It checks if there is an active database connection (s.DB) and attempts to close it.
// Logs an error message if there was an issue closing the connection.
func (s *PSQLClient) Close() {
	if s.DB != nil {
		if err := s.DB.Close(); err != nil {
			log.Errorf("Error closing connection to PostgreSQL: %v", err)
		} else {
			log.Println("Connection to PostgreSQL closed")
		}
	}
}

// SqlMigrate runs database migrations for PostgreSQL using the provided database connection.
//
// It initializes the migration driver with the current database instance, and then
// applies all available migrations from the specified migrations directory ("migrations").
//
// Returns an error if the database connection (s.DB) is nil, if there is an issue
// initializing the migration driver, or if there are errors during the migration process.
func (s *PSQLClient) SqlMigrate() error {
	const op = "storage.postgres.SqlMigrate()"

	if s.DB == nil {
		return fmt.Errorf("%s DB is nil", op)
	}

	driver, err := postgres.WithInstance(s.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
