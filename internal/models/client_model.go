package models

import "time"

type Client struct {
	ID          int64
	ClientName  string
	Version     int
	Image       string
	CPU         string
	Memory      string
	Priority    float64
	NeedRestart bool
	SpawnedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
