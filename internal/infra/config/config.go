package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once sync.Once
)

type config struct {
	Database  DatabaseConfig
	Server    ServerConfig
	Log       LogConfig
	Telemetry TelemetryConfig
}

type ConfigInterface interface {
	GetDatabaseConfig() DatabaseConfig
	GetServerConfig() ServerConfig
	GetLogConfig() LogConfig
	GetTelemetryConfig() TelemetryConfig
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SslMode         string
	MaxConnections  int
	MinConnections  int
	ConnMaxLifetime int
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	Environment  string
}

type LogConfig struct {
	Level       string
	Path        string
	ServiceName string
	Environment string
}

type TelemetryConfig struct {
	ServiceName              string
	ServiceVersion           string
	OtelExporterOtlpEndpoint string
	OtelExporterOtlpInsecure bool
}

func New() ConfigInterface {
	var cfg *config
	once.Do(func() {
		_ = godotenv.Load()

		cfg = &config{
			Database: DatabaseConfig{
				Host:            getEnv("DB_HOST", "localhost"),
				Port:            getEnv("DB_PORT", "5432"),
				User:            getEnv("DB_USER", "user"),
				Password:        getEnv("DB_PASSWORD", ""),
				Name:            getEnv("DB_NAME", "db"),
				SslMode:         getEnv("DB_SSL", "disable"),
				MaxConnections:  getEnvInt("DB_MAX_CONNECTIONS", 20),
				MinConnections:  getEnvInt("DB_MIN_CONNECTIONS", 1),
				ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 300),
			},
			Server: ServerConfig{
				Port:         getEnv("SERVER_PORT", "8000"),
				ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 15),
				WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 15),
				IdleTimeout:  getEnvInt("SERVER_IDLE_TIMEOUT", 60),
				Environment:  getEnv("ENVIRONMENT", "development"),
			},
			Log: LogConfig{
				Level:       getEnv("LOG_LEVEL", "INFO"),
				Path:        getEnv("LOG_PATH", "./logs/app.log"),
				ServiceName: getEnv("SERVICE_NAME", "auth-service"),
				Environment: getEnv("ENVIRONMENT", "development"),
			},
			Telemetry: TelemetryConfig{
				ServiceName:              getEnv("SERVICE_NAME", "auth-service"),
				ServiceVersion:           getEnv("OTEL_SERVICE_VERSION", "0.0.1"),
				OtelExporterOtlpEndpoint: getEnv("OTEL_EXPORTER_ENDPOINT", "http://collector:4317"),
				OtelExporterOtlpInsecure: true,
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

func (c *config) GetTelemetryConfig() TelemetryConfig {
	return c.Telemetry
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue != "" {
			return defaultValue
		}
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}
