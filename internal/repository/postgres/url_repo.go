package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"url-shortener/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresURLRepository struct {
	db *sql.DB
}

// NewPostgresURLRepository creates a new postgres repository
func NewPostgresURLRepository(db *sql.DB) domain.URLRepository {
	return &postgresURLRepository{
		db: db,
	}
}

func (r *postgresURLRepository) Store(ctx context.Context, u *domain.URL) error {
	query := `
		INSERT INTO urls (short_code, original_url, click_count, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, u.ShortCode, u.OriginalURL, u.ClickCount, time.Now()).Scan(&u.ID)
	return err
}

func (r *postgresURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	query := `
		SELECT id, short_code, original_url, click_count, created_at
		FROM urls
		WHERE short_code = $1
	`
	var u domain.URL
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ClickCount, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("url not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresURLRepository) IncrementClick(ctx context.Context, shortCode string) error {
	query := `
		UPDATE urls
		SET click_count = click_count + 1
		WHERE short_code = $1
	`
	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("url not found")
	}
	return nil
}

func (r *postgresURLRepository) GetStats(ctx context.Context, shortCode string) (*domain.URLStats, error) {
	query := `
		SELECT short_code, click_count
		FROM urls
		WHERE short_code = $1
	`
	var s domain.URLStats
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&s.ShortCode, &s.ClickCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("url not found")
		}
		return nil, err
	}
	return &s, nil
}
