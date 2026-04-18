package configs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Params   string
}

func LoadDatabaseConfig() (DatabaseConfig, error) {
	_ = godotenv.Load()

	cfg := DatabaseConfig{
		Driver:   envOrDefault("DB_DRIVER", "mysql"),
		Host:     envOrDefault("DB_HOST", "127.0.0.1"),
		Port:     envOrDefault("DB_PORT", "3306"),
		User:     envOrDefault("DB_USER", "root"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     envOrDefault("DB_NAME", "govibe"),
		Params:   envOrDefault("DB_PARAMS", "parseTime=true&loc=Local"),
	}

	if strings.TrimSpace(cfg.Driver) == "" {
		return DatabaseConfig{}, errors.New("DB_DRIVER is required")
	}
	if cfg.Driver != "mysql" {
		return DatabaseConfig{}, fmt.Errorf("unsupported DB_DRIVER %q (supported: mysql)", cfg.Driver)
	}
	if strings.TrimSpace(cfg.Host) == "" {
		return DatabaseConfig{}, errors.New("DB_HOST is required")
	}
	if strings.TrimSpace(cfg.Port) == "" {
		return DatabaseConfig{}, errors.New("DB_PORT is required")
	}
	if strings.TrimSpace(cfg.User) == "" {
		return DatabaseConfig{}, errors.New("DB_USER is required")
	}
	if strings.TrimSpace(cfg.Name) == "" {
		return DatabaseConfig{}, errors.New("DB_NAME is required")
	}

	return cfg, nil
}

func (c DatabaseConfig) MySQLDSN() string {
	params := strings.TrimPrefix(strings.TrimSpace(c.Params), "?")
	if params == "" {
		params = "parseTime=true&loc=Local"
	}

	// Keep it simple: user/password may contain special chars; if you need full escaping,
	// set DB_DSN and use that directly in your DB setup.
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.User, c.Password, c.Host, c.Port, c.Name, params)
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
