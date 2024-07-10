package repository

import "database/sql"

type AlgorithmRepository interface {
}

type algorithmRepository struct {
	db *sql.DB
}

func NewAlgorithmRepository(db *sql.DB) AlgorithmRepository {
	return &algorithmRepository{db: db}
}
