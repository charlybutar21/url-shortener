package usecase

import (
	"context"
	"testing"

	"url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Store(ctx context.Context, url *domain.URL) error {
	args := m.Called(ctx, url)
	if url != nil && args.Error(0) == nil {
		url.ID = 1 // simulate DB auto-increment
	}
	return args.Error(0)
}

func (m *MockURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.URL), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockURLRepository) IncrementClick(ctx context.Context, shortCode string) error {
	args := m.Called(ctx, shortCode)
	return args.Error(0)
}

func (m *MockURLRepository) GetStats(ctx context.Context, shortCode string) (*domain.URLStats, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.URLStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestShortenURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	uc := NewURLUsecase(mockRepo)

	// Test case: Success
	mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.URL")).Return(nil).Once()
	
	result, err := uc.ShortenURL(context.Background(), "https://example.com")
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.ShortCode, 7) // Ensure it's 7 chars
	assert.Equal(t, "https://example.com", result.OriginalURL)

	// Test case: Invalid URL
	_, err = uc.ShortenURL(context.Background(), "invalid-url")
	assert.Error(t, err)

	// Test case: Empty URL
	_, err = uc.ShortenURL(context.Background(), "")
	assert.Error(t, err)
}

func TestGetOriginalURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	uc := NewURLUsecase(mockRepo)

	mockURL := &domain.URL{
		ShortCode:   "abcdefg",
		OriginalURL: "https://example.com",
	}

	mockRepo.On("GetByShortCode", mock.Anything, "abcdefg").Return(mockURL, nil)
	mockRepo.On("IncrementClick", mock.Anything, "abcdefg").Return(nil)

	urlStr, err := uc.GetOriginalURL(context.Background(), "abcdefg")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", urlStr)
}
