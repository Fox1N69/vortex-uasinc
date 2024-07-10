package main

import (
	"test-task/infra"
	"test-task/internal/api"
	"test-task/pkg/util/logger"
)

// User algoritm sync
func main() {
	// Init config
	i := infra.New("config/config.json")
	// Set project mod
	mode := i.SetMode()

	//init logger
	logger.Init(mode)
	log := logger.GetLogger()

	//Connect to database and migration
	i.PSQLClient()
	log.Info("Connect to PSQLClient")
	i.RunSQLMigrations()
	log.Info("Sql migrations compelet")

	// Start api server
	api.NewServer(i).Run()
}
