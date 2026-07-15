package usecase

import (
	"context"
	"errors"
	"net/url"
	
	"url-shortener/internal/domain"
	"url-shortener/pkg/base62"
)

type urlUsecase struct {
	urlRepo domain.URLRepository
}

// NewURLUsecase creates a new url usecase
func NewURLUsecase(repo domain.URLRepository) domain.URLUsecase {
	return &urlUsecase{
		urlRepo: repo,
	}
}

func (u *urlUsecase) ShortenURL(ctx context.Context, originalURL string) (*domain.URL, error) {
	if originalURL == "" {
		return nil, errors.New("url cannot be empty")
	}
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return nil, errors.New("invalid url format")
	}

	id, err := u.urlRepo.GenerateID(ctx)
	if err != nil {
		return nil, err
	}

	shortCode := base62.Encode(id)

	newURL := &domain.URL{
		ID:          id,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
	}

	err = u.urlRepo.Store(ctx, newURL)
	if err != nil {
		return nil, err
	}

	return newURL, nil
}

func (u *urlUsecase) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	if shortCode == "" {
		return "", errors.New("short code cannot be empty")
	}

	urlObj, err := u.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// Increment click asynchronously
	go func() {
		_ = u.urlRepo.IncrementClick(context.Background(), shortCode)
	}()

	return urlObj.OriginalURL, nil
}

func (u *urlUsecase) GetURLStats(ctx context.Context, shortCode string) (*domain.URLStats, error) {
	if shortCode == "" {
		return nil, errors.New("short code cannot be empty")
	}
	return u.urlRepo.GetStats(ctx, shortCode)
}
