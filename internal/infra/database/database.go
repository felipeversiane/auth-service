package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	once     sync.Once
	instance *database
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
		logger.Info("Initializing database connection...")

		dsn := getConnectionString(config)

		poolConfig, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse pool config: %w", parseErr)
			logger.Error("Error parsing pool config", zap.Error(err))
			return
		}

		poolConfig.MaxConns = int32(config.MaxConnections)
		poolConfig.MinConns = int32(config.MinConnections)
		poolConfig.MaxConnLifetime = time.Duration(config.ConnMaxLifetime) * time.Second
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

		logger.Info("Creating database connection pool...")

		pool, connErr := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if connErr != nil {
			err = fmt.Errorf("failed to create connection pool: %w", connErr)
			logger.Error("Error creating connection pool", zap.Error(err))
			return
		}

		instance = &database{
			db:     pool,
			config: config,
		}

		logger.Info("Attempting to ping database...")

		if err := instance.Ping(context.Background()); err != nil {
			instance.Close()
			err = fmt.Errorf("failed to connect to database: %w", err)
			logger.Error("Error connecting to database", zap.Error(err))
		} else {
			logger.Info("Database connection established successfully")
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
		logger.Warn("Database ping failed", zap.Error(err))
	}
	return err
}

func (d *database) Close() {
	if d.db != nil {
		d.db.Close()
		logger.Info("Database connection closed")
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
