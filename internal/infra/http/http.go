package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/database"
	"github.com/gin-gonic/gin"
)

type httpServer struct {
	router *gin.Engine
	srv    *http.Server
	config config.HttpServerConfig
	db     database.DatabaseInterface
}

type HttpServerInterface interface {
	Start() error
	Shutdown(ctx context.Context) error
	InitRoutes()
}

func New(config config.HttpServerConfig, db database.DatabaseInterface) HttpServerInterface {
	if config.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	server := &httpServer{
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

func (s *httpServer) InitRoutes() {
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

func (s *httpServer) Start() error {
	slog.Info("starting HTTP server", slog.String("port", s.config.Port))

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		err = fmt.Errorf("failed to start server: %w", err)
		slog.Error("server not started", "error", err)
		return err
	}

	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	slog.Info("initiating graceful shutdown")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("error during server shutdown: %w", err)
		slog.Error("server shutdown unsuccessfully", "error", err)
		return err
	}

	slog.Info("server shutdown completed successfully")
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
