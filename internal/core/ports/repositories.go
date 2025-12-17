package ports

import (
	"context"

	"github.com/elghazx/go-shorturl/internal/core/domain"
)

type URLRepository interface {
	Save(ctx context.Context, url *domain.URL) error
	GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error)
	IncrementClicks(ctx context.Context, shortCode string) error
	GetTopURLs(ctx context.Context, limit int) ([]domain.URL, error)
}

type CacheService interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
}
