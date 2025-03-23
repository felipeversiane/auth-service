package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once     sync.Once
	instance *database
)

const (
	maxConns = 20
	minConns = 2
)

type database struct {
	db     *pgxpool.Pool
	config config.DatabaseConfig
}

type DatabaseInterface interface {
	GetDB() *pgxpool.Pool
	Ping(ctx context.Context) error
	Close()
}

func New(config config.DatabaseConfig) (DatabaseInterface, error) {
	var err error
	once.Do(func() {
		slog.Info("initializing database connection...")

		dsn := getConnectionString(config)

		poolConfig, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse pool config: %w", parseErr)
			slog.Error("error parsing pool config", "error", err)
			return
		}

		poolConfig.MaxConns = maxConns
		poolConfig.MinConns = minConns

		slog.Info("creating database connection pool...")

		pool, connErr := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if connErr != nil {
			err = fmt.Errorf("failed to create connection pool: %w", connErr)
			slog.Error("error creating connection pool", "error", err)
			return
		}

		instance = &database{
			db:     pool,
			config: config,
		}

		slog.Info("attempting to ping database...")

		if err := instance.Ping(context.Background()); err != nil {
			instance.Close()
			err = fmt.Errorf("failed to connect to database: %w", err)
			slog.Error("error connecting to database", "error", err)
		} else {
			slog.Info("database connection established successfully")
		}
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (d *database) Ping(ctx context.Context) error {
	err := d.db.Ping(ctx)
	if err != nil {
		slog.Warn("database ping failed", "error", err)
	}
	return err
}

func (d *database) Close() {
	if d.db != nil {
		d.db.Close()
		slog.Info("database connection closed")
	}
}

func (d *database) GetDB() *pgxpool.Pool {
	return d.db
}

func getConnectionString(config config.DatabaseConfig) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s",
		config.User,
		config.Password,
		config.Name,
		config.Port,
		config.Host,
		config.SslMode,
	)
}
