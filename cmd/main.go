package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Woolfer0097/InMemoryLinks"
	appRouter "github.com/Woolfer0097/UrlShortener/internal"
	"github.com/Woolfer0097/UrlShortener/internal/config"
	"github.com/Woolfer0097/UrlShortener/internal/repository"
	"github.com/Woolfer0097/UrlShortener/internal/repository/models"
	postgresrepo "github.com/Woolfer0097/UrlShortener/internal/repository/postgres"
	"github.com/Woolfer0097/UrlShortener/pkg/logger"
)

func main() {
	cfg := config.GetConfig()
	logger.Init(true)

	var urlRepo repository.UrlRepository
	if cfg.Storage.Type == config.StorageInMemory {
		urlRepo = repository.NewInMemoryLinksAdapter(inmemorylinks.NewUrlRepository())
	} else {
		dsn := config.GetPostgresDSN()
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
			log.Fatalf("failed to enable uuid extension: %v", err)
		}
		if err := db.AutoMigrate(&models.Url{}); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
		urlRepo = postgresrepo.NewUrlRepository(db)
	}

	r := chi.NewRouter()
	appRouter.InitRoutes(r, urlRepo)

	addr := cfg.API.AppHost + ":" + cfg.API.AppPort
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
