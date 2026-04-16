package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TASKPORT     string
	DATABASE_URL string
}

func NewСonfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{TASKPORT: os.Getenv("TASKPORT"), DATABASE_URL: os.Getenv("DATABASE_URL")}, nil
}
