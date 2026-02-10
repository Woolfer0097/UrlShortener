package postgres

import (
	"context"

	"github.com/Woolfer0097/UrlShortener/internal/repository/models"

	"gorm.io/gorm"
)

type UrlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) *UrlRepository {
	return &UrlRepository{db: db}
}

func (r *UrlRepository) Create(ctx context.Context, url *models.Url) error {
	return r.db.WithContext(ctx).Create(url).Error
}

func (r *UrlRepository) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	var url models.Url
	err := r.db.WithContext(ctx).
		First(&url, "url_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *UrlRepository) GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error) {
	var url models.Url
	err := r.db.WithContext(ctx).
		First(&url, "original_url = ?", originalUrl).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}
