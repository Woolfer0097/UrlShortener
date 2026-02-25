package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/Woolfer0097/UrlShortener/internal/config"
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

func isValidURL(url string) bool {
	return len(url) <= 2048
}

func (h *ShortenerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req schemas.ShortenUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OriginalUrl == "" || !isValidURL(req.OriginalUrl) {
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

	maxCreateAttempts := config.MaxCreateAttempts
	baseURL := strings.TrimSuffix(h.BaseURL, "/")
	for attempt := 1; attempt <= maxCreateAttempts; attempt++ {
		code := shortener.GenerateShortUrl(originalUrl)
		existingByCode, err := h.UrlRepo.GetByCode(r.Context(), code)
		if err == nil && existingByCode != nil {
			if existingByCode.OriginalUrl == originalUrl {
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(schemas.ShortenUrlResponse{ShortUrl: baseURL + "/" + code})
				if err != nil {
					http.Error(w, "Shorten URL: Bad response format", http.StatusInternalServerError)
					return
				}
				return
			}
			continue
		}
		url := &models.Url{UrlCode: code, OriginalUrl: originalUrl}
		if err := h.UrlRepo.Create(r.Context(), url); err != nil {
			if attempt == maxCreateAttempts {
				http.Error(w, "Max attempts amount used and code hasn't been generated till now", http.StatusInternalServerError)
				return
			}
			continue
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(schemas.ShortenUrlResponse{ShortUrl: baseURL + "/" + code})
		if err != nil {
			http.Error(w, "Shorten URL: Bad response format", http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func (h *ShortenerHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{10}$`)
	if !re.MatchString(code) {
		http.Error(w, "Invalid code format", http.StatusBadRequest)
		return
	}
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
	err = json.NewEncoder(w).Encode(
		schemas.GetUrlResponse{
			OriginalUrl: url.OriginalUrl,
		},
	)
	if err != nil {
		http.Error(w, "Get URL: Bad response format", http.StatusInternalServerError)
		return
	}
}
