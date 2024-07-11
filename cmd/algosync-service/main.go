package main

import (
	"net/http"
	_ "net/http/pprof"
	_ "test-task/cmd/algosync-service/docs"
	"test-task/infra"
	"test-task/internal/api"
)

// @title AlgorithmSync service
// @description сервис для синхронизации пользовательских алгоритмов
// @version 1.0

// @host localhost:4000
// @basePath /api
func main() {

	// Init config
	i := infra.New("config/config.json")
	// Set project mod
	i.SetMode()

	// Get custom logrus logger
	log := i.GetLogger()

	// Start pprof server
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//Connect to database and migration
	i.PSQLClient()
	log.Info("Connect to PSQLClient")
	i.RunSQLMigrations()
	log.Info("Sql migrations compelet")

	// Start api server
	api.NewServer(i).Run()
}
