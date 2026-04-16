package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TASKPORT     int
	DATABASE_URL string
}

func NewСonfig() (*Config, error) {
	godotenv.Load()
	TaskPort, err := strconv.Atoi(os.Getenv("TASKPORT"))
	if err != nil {
		return nil, err
	}
	return &Config{TASKPORT: TaskPort, DATABASE_URL: os.Getenv("DATABASE_URL")}, nil
}
