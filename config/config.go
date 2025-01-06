package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

func LoadTestConfig() (*DatabaseConfig, error) {
	if err := godotenv.Load(".env.test"); err != nil {
		return nil, fmt.Errorf("error loading .env.test file: %w", err)
	}

	cfg := &DatabaseConfig{
		Postgres: PostgresConfig{
			Host:     os.Getenv("TEST_DB_HOST"),
			Port:     os.Getenv("TEST_DB_PORT"),
			User:     os.Getenv("TEST_DB_USER"),
			Password: os.Getenv("TEST_DB_PASS"),
			DBName:   os.Getenv("TEST_DB_NAME"),
		},
	}

	return cfg, nil
}
