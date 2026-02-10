package repository

import (
	"context"

	"github.com/Woolfer0097/UrlShortener/internal/repository/models"
)

type UrlRepository interface {
	Create(ctx context.Context, url *models.Url) error
	GetByCode(ctx context.Context, code string) (*models.Url, error)
	GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error)
}
