package domain

import (
	"context"
	"time"
)

// URL represents the url entity
type URL struct {
	ID          uint64    `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ClickCount  uint64    `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// URLStats represents the statistics of a short url
type URLStats struct {
	ShortCode  string `json:"short_code"`
	ClickCount uint64 `json:"click_count"`
}

// URLRepository represent the repository contract
type URLRepository interface {
	Store(ctx context.Context, url *URL) error
	GetByShortCode(ctx context.Context, shortCode string) (*URL, error)
	IncrementClick(ctx context.Context, shortCode string) error
	GetStats(ctx context.Context, shortCode string) (*URLStats, error)
}

// URLUsecase represent the usecase contract
type URLUsecase interface {
	ShortenURL(ctx context.Context, originalURL string) (*URL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	GetURLStats(ctx context.Context, shortCode string) (*URLStats, error)
}
