package repository

import (
	"context"
	"sync"

	inmemorylinks "github.com/Woolfer0097/InMemoryLinks"
	"github.com/Woolfer0097/UrlShortener/internal/repository/models"
	"gorm.io/gorm"
)

type urlRepoAdapter struct {
	mu            sync.RWMutex
	inner         *inmemorylinks.UrlRepository
	byOriginalUrl map[string]string
}

func NewInMemoryLinksAdapter(inner *inmemorylinks.UrlRepository) UrlRepository {
	return &urlRepoAdapter{inner: inner, byOriginalUrl: make(map[string]string)}
}

func (a *urlRepoAdapter) Create(ctx context.Context, url *models.Url) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	u := &inmemorylinks.Url{
		ID:          url.ID,
		UrlCode:     url.UrlCode,
		OriginalUrl: url.OriginalUrl,
	}

	if err := a.inner.Create(ctx, u); err != nil {
		return err
	}

	a.byOriginalUrl[url.OriginalUrl] = url.UrlCode
	return nil
}

func (a *urlRepoAdapter) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	u, err := a.inner.GetByCode(ctx, code)
	if err != nil {
		if err == inmemorylinks.ErrNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &models.Url{
		ID:          u.ID,
		UrlCode:     u.UrlCode,
		OriginalUrl: u.OriginalUrl,
	}, nil
}

func (a *urlRepoAdapter) GetByOriginalUrl(ctx context.Context, originalUrl string) (*models.Url, error) {
	a.mu.RLock()
	code, ok := a.byOriginalUrl[originalUrl]
	a.mu.RUnlock()
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return a.GetByCode(ctx, code)
}
