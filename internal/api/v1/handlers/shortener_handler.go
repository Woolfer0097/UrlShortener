package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Woolfer0097/UrlShortener/internal/repository"
	"github.com/Woolfer0097/UrlShortener/internal/repository/models"
	"github.com/Woolfer0097/UrlShortener/internal/repository/schemas"
	shortener "github.com/Woolfer0097/UrlShortener/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ShortenerHandler struct {
	UrlRepo repository.UrlRepository
	BaseURL string
}

func (h *ShortenerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req schemas.ShortenUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OriginalUrl == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	originalUrl := req.OriginalUrl

	existing, err := h.UrlRepo.GetByOriginalUrl(r.Context(), originalUrl)
	if err == nil && existing != nil {
		baseURL := strings.TrimSuffix(h.BaseURL, "/")
		shortURL := baseURL + "/" + existing.UrlCode
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schemas.ShortenUrlResponse{ShortUrl: shortURL})
		return
	}

	const maxCreateAttempts = 5
	baseURL := strings.TrimSuffix(h.BaseURL, "/")
	for attempt := 0; attempt < maxCreateAttempts; attempt++ {
		code := shortener.GenerateShortUrl(originalUrl)
		existingByCode, err := h.UrlRepo.GetByCode(r.Context(), code)
		if err == nil && existingByCode != nil {
			if existingByCode.OriginalUrl == originalUrl {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(schemas.ShortenUrlResponse{ShortUrl: baseURL + "/" + code})
				return
			}
			continue
		}
		url := &models.Url{UrlCode: code, OriginalUrl: originalUrl}
		if err := h.UrlRepo.Create(r.Context(), url); err != nil {
			if attempt == maxCreateAttempts-1 {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			continue
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schemas.ShortenUrlResponse{ShortUrl: baseURL + "/" + code})
		return
	}
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func (h *ShortenerHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	url, err := h.UrlRepo.GetByCode(r.Context(), code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		schemas.GetUrlResponse{
			OriginalUrl: url.OriginalUrl,
		},
	)
}
