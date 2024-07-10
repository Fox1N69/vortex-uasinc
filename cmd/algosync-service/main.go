package main

import (
	"test-task/infra"
	"test-task/internal/api"
)

// User algoritm sync
func main() {
	// Init config
	i := infra.New("config/config.json")
	// Set project mod
	i.SetMode()

	//init logger
	log := i.GetLogger()

	//Connect to database and migration
	i.PSQLClient()
	log.Info("Connect to PSQLClient")
	i.RunSQLMigrations()
	log.Info("Sql migrations compelet")

	// Start api server
	api.NewServer(i).Run()
}
