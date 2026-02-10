package routes

import (
	"github.com/Woolfer0097/UrlShortener/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterShortenerRoutes(r chi.Router, h *handlers.ShortenerHandler) {
	r.Post("/shorten", h.Create)
	r.Get("/{code}", h.Get)
}
