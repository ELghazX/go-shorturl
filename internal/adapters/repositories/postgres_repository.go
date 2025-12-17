package repositories

import (
	"context"
	"database/sql"

	"github.com/elghazx/go-shorturl/internal/core/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Save(ctx context.Context, url *domain.URL) error {
	query := `INSERT INTO urls (short_code, long_url) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, url.ShortCode, url.LongURL)
	return err
}

func (r *PostgresRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	query := `SELECT id, short_code, long_url, clicks, created_at FROM urls WHERE short_code = $1`

	url := &domain.URL{}
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(
		&url.ID, &url.ShortCode, &url.LongURL, &url.Clicks, &url.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (r *PostgresRepository) IncrementClicks(ctx context.Context, shortCode string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = $1`
	_, err := r.db.ExecContext(ctx, query, shortCode)
	return err
}

func (r *PostgresRepository) GetTopURLs(ctx context.Context, limit int) ([]domain.URL, error) {
	query := `SELECT id, short_code, long_url, clicks, created_at FROM urls ORDER BY clicks DESC LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []domain.URL

	for rows.Next() {
		var url domain.URL
		err := rows.Scan(&url.ID, &url.ShortCode, &url.LongURL, &url.Clicks, &url.CreatedAt)
		if err != nil {
			continue
		}
		urls = append(urls, url)
	}
	return urls, nil
}
