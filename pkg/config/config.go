package config

import (
	"fmt"
	"os"

	"go1f/pkg/logger"

	"github.com/joho/godotenv"
)

var log = logger.GetLogger()

type Config struct {
	ToDoPort   string
	ToDoDBFile string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	var cfg Config
	cfg.ToDoPort = os.Getenv("TODO_PORT")

	if cfg.ToDoPort == "" {
		cfg.ToDoPort = "7540"
	}

	cfg.ToDoDBFile = os.Getenv("TODO_DBFILE")

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	log.Info("Environment variables loaded successfully")
	return &cfg, nil
}

func validateConfig(cfg *Config) error {

	if cfg.ToDoDBFile == "" {
		return fmt.Errorf("missing required field 'TODO_DBFILE'")
	}

	return nil
}
