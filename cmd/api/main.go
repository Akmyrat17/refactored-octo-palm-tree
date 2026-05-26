package main

import (
	"context"
	"fmt"
	"log"

	"github.com/boilerplate/internal/config"
	"github.com/boilerplate/internal/platform/database"
	"github.com/boilerplate/internal/server"
	"github.com/boilerplate/pkg/logger"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(context.Background(), cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log := logger.NewConsoleLogger()

	e := server.New(cfg, db, log)

	fmt.Printf("Starting server on port %d\n", cfg.Server.Port)
	if err := server.Start(e, cfg.Server.Port); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
