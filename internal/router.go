package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Woolfer0097/UrlShortener/internal/api/v1/handlers"
	"github.com/Woolfer0097/UrlShortener/internal/config"
	"github.com/Woolfer0097/UrlShortener/internal/api/v1/routes"
	"github.com/Woolfer0097/UrlShortener/internal/repository"
)

func InitRoutes(r *chi.Mux, urlRepo repository.UrlRepository) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	urlHandler := &handlers.ShortenerHandler{
		UrlRepo: urlRepo,
		BaseURL: config.GetConfig().API.BaseURL,
	}
	routes.RegisterShortenerRoutes(r, urlHandler)
}
