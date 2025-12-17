package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/elghazx/go-shorturl/internal/core/domain"
	"github.com/elghazx/go-shorturl/internal/core/ports"
)

type URLService struct {
	repo  ports.URLRepository
	cache ports.CacheService
}

func NewURLService(repo ports.URLRepository, cache ports.CacheService) *URLService {
	return &URLService{
		repo:  repo,
		cache: cache,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, longURL string) (*domain.URL, error) {
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "https://" + longURL
	}

	url := &domain.URL{
		ShortCode: s.generateShortCode(),
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, err
	}
	s.cache.Set(ctx, url.ShortCode, url.LongURL)
	return url, nil
}

func (s *URLService) generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}

func (s *URLService) GetStats(ctx context.Context) ([]domain.URL, error) {
	return s.repo.GetTopURLs(ctx, 10)
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	// cache
	if longURL, err := s.cache.Get(ctx, shortCode); err == nil {
		go s.repo.IncrementClicks(context.Background(), shortCode)
		return longURL, nil
	}

	// use db instead
	url, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	s.cache.Set(ctx, shortCode, url.LongURL)
	go s.repo.IncrementClicks(context.Background(), shortCode)
	return url.LongURL, nil
}
