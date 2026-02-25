package main

import (
	"context"
	"database/sql"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	inmemorylinks "github.com/Woolfer0097/InMemoryLinks"
	appRouter "github.com/Woolfer0097/UrlShortener/internal"
	"github.com/Woolfer0097/UrlShortener/internal/config"
	"github.com/Woolfer0097/UrlShortener/internal/repository"
	postgresrepo "github.com/Woolfer0097/UrlShortener/internal/repository/postgres"
	"github.com/Woolfer0097/UrlShortener/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.GetConfig()
	logger.Init(true)

	var sqlDB *sql.DB

	var urlRepo repository.UrlRepository
	if cfg.Storage.Type == config.StorageInMemory {
		urlRepo = repository.NewInMemoryLinksAdapter(inmemorylinks.NewUrlRepository())
	} else {
		dsn := config.GetPostgresDSN()
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Log.Error("failed to connect to database", zap.Error(err))
		}

		// for graceful shutdown
		sqlDB, err = db.DB()
		if err != nil {
			logger.Log.Error("failed to get sql.DB", zap.Error(err))
		}

		goose.SetDialect("postgres")

		// Run all migrations in migrations folder
		if err := goose.Up(sqlDB, "migrations"); err != nil {
			logger.Log.Error("failed to apply migrations", zap.Error(err))
		}

		logger.Log.Info("Migrations applied successfully!")

		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
			logger.Log.Error("failed to enable uuid extension", zap.Error(err))
		}
		// if err := db.AutoMigrate(&models.Url{}); err != nil {
		// 	logger.Log.Error("failed to migrate database", zap.Error(err))
		// }
		urlRepo = postgresrepo.NewUrlRepository(db)
	}

	r := chi.NewRouter()
	appRouter.InitRoutes(r, urlRepo)

	addr := cfg.API.AppHost + ":" + cfg.API.AppPort

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server error", zap.Error(err))
			return
		}
		logger.Log.Info("server started", zap.String("addr", addr), zap.String("host", cfg.API.AppHost), zap.String("port", cfg.API.AppPort))
	}()

	<-ctx.Done()
	logger.Log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("graceful shutdown failed", zap.Error(err))
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			logger.Log.Error("error closing database", zap.Error(err))
		}
	}

	logger.Log.Info("server stopped gracefully")
}
