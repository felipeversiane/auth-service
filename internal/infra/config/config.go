package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once sync.Once
)

type config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Log      LogConfig
}

type ConfigInterface interface {
	GetDatabaseConfig() DatabaseConfig
	GetServerConfig() ServerConfig
	GetLogConfig() LogConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
}

type ServerConfig struct {
	Port string
}

type LogConfig struct {
	Level string
}

func New() ConfigInterface {
	var cfg *config
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			slog.Warn("no .env file found")
		}

		cfg = &config{
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "user"),
				Password: getEnv("DB_PASSWORD", ""),
				Name:     getEnv("DB_NAME", "db"),
				SslMode:  getEnv("DB_SSL", "disable"),
			},
			Server: ServerConfig{
				Port: getEnv("SERVER_PORT", "8000"),
			},
			Log: LogConfig{
				Level: getEnv("LOG_LEVEL", "INFO"),
			},
		}
	})

	return cfg
}

func (c *config) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}

func (c *config) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *config) GetLogConfig() LogConfig {
	return c.Log
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue != "" {
			return defaultValue
		}
		slog.Error("missing required environment variable", "key", key)
	}
	return value
}
