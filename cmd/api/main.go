package main

import (
	"log"

	"github.com/personal/task-management/cmd/api/wire"
)

// @title Task Management API
// @version 1.0
// @description Task Management API
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	app, cleanup, err := wire.NewWire()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
