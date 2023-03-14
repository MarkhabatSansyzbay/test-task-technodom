package main

import (
	"context"
	"fmt"
	"log"

	"task/internal/config"
	"task/internal/repository"
)

const configPath = "./config/config.toml"

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

	fmt.Println(db)
}
