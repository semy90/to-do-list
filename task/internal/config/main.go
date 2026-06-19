package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	DATABASE_URL string
}

func NewСonfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{PORT: os.Getenv("PORT"), DATABASE_URL: os.Getenv("DATABASE_URL")}, nil
}
