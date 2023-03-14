package main

import (
	"context"
	"fmt"
	"log"

	"task/internal/config"
	"task/internal/repository"
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
	fmt.Println(service.SaveDataset(datasetFile))
	fmt.Println(db)
}
