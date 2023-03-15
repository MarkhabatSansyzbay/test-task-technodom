package main

import (
	"context"
	"fmt"
	"log"

	"task/internal/config"
	"task/internal/delivery"
	"task/internal/repository"
	"task/internal/server"
	"task/internal/service"
)

const (
	configPath  = "./config/config.toml"
	datasetFile = "./links.json"
)

func main() {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("config: %s\n", err)
	}

	ctx := context.Background()
	db, err := repository.PostgreSQLDB(cfg.Dsn, ctx)
	if err != nil {
		log.Fatalf("database initialization error: %s\n", err)
	}

	repo := repository.NewRedirecter(db)
	service := service.NewRedirecter(repo)
	handler := delivery.NewHandler(service)
	server := new(server.Server)

	if err := service.SaveDataset(datasetFile); err != nil {
		log.Printf("error inserting dataset to db: %s", err)
	}

	fmt.Println("Starting server...")

	if err := server.Run(cfg.Port, handler.InitRoutes()); err != nil {
		log.Fatalf("error while running the server: %s", err.Error())
	}
}
