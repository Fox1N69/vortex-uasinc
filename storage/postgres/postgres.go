package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type PSQLClient struct {
	DB *sql.DB
}

func NewPSQLClient() *PSQLClient {
	return &PSQLClient{}
}

// Connect establishes a connection to the PostgreSQL database
func (s *PSQLClient) Connect(user, password, host, port, dbname string) error {
	const op = "storage.postgres.Connect()"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("%s %w", op, err)
	}

	s.DB = db
	log.Println("Connect to Postgres")
	return nil
}

// Close terminates the connection to the database
func (s *PSQLClient) Close() {
	if s.DB != nil {
		if err := s.DB.Close(); err != nil {
			log.Errorf("Error closing connection to PostgreSQL: %v", err)
		} else {
			log.Println("Connection to PostgreSQL closed")
		}
	}
}

// SqlMigrate performs database schema migrations using Golang Migrate
func (s *PSQLClient) SqlMigrate() error {
	const op = "storage.postgres.SqlMigrate()"

	driver, err := postgres.WithInstance(s.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s %w", op, err)
	}

	log.Println("SQL migrations completed")
	return nil
}
