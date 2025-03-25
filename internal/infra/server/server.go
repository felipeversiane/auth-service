package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/database"
	"github.com/felipeversiane/auth-service/internal/infra/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	router *gin.Engine
	srv    *http.Server
	config config.ServerConfig
	db     database.DatabaseInterface
}

type ServerInterface interface {
	Start() error
	Shutdown(ctx context.Context) error
	InitRoutes()
}

func New(config config.ServerConfig, db database.DatabaseInterface) ServerInterface {
	if config.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	server := &server{
		router: router,
		srv: &http.Server{
			Addr:         ":" + config.Port,
			Handler:      router,
			ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
		},
		config: config,
		db:     db,
	}

	return server
}

func (s *server) InitRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"status":    "up",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})
	}
}

func (s *server) Start() error {
	logger.Info("Starting HTTP server", zap.String("port", s.config.Port))

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	logger.Info("Initiating graceful shutdown")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("error during server shutdown: %w", err)
	}

	logger.Info("Server shutdown completed successfully")
	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
