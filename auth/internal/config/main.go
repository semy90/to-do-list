package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	DATABASE_URL string
	REDIS_URL    string
	SALT         string
}

func NewСonfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{REDIS_URL: os.Getenv("REDIS_URL"), PORT: os.Getenv("PORT"), DATABASE_URL: os.Getenv("DATABASE_URL"), SALT: os.Getenv("SALT")}, nil
}
