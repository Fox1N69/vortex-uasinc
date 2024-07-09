package main

import (
	"test-task/infra"
	"test-task/internal/api"
)

// User algoritm sync
func main() {
	// Init config
	i := infra.New("/config/config.json")
	// Start api server
	api.NewServer(i).Run()
}
