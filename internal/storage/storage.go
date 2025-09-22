package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v8"
)

type Store struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

type URLRecord struct {
	ID         int       `db:"id"`
	Code       string    `db:"code"`
	Original   string    `db:"original_url"`
	CreatedAt  time.Time `db:"created_at"`
	ExpiresAt  *time.Time `db:"expires_at"`
	ClicksCount int64    `db:"clicks_count"`
}

func NewStore(db *sqlx.DB, r *redis.Client) *Store {
	return &Store{DB: db, Redis: r}
}

func (s *Store) CreateURL(ctx context.Context, code, original string, expiresAt *time.Time) (int, error) {
	var id int
	query := `INSERT INTO urls (code, original_url, expires_at) VALUES ($1,$2,$3) RETURNING id`
	err := s.DB.QueryRowxContext(ctx, query, code, original, expiresAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) GetURLByCode(ctx context.Context, code string) (*URLRecord, error) {
	// Optionally check Redis cache first
	const cacheKeyPrefix = "shortner:url:"
	cacheKey := cacheKeyPrefix + code
	if s.Redis != nil {
		val, err := s.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			// cache hit -> return minimal record with Original
			return &URLRecord{Code: code, Original: val}, nil
		}
		// ignore redis errors and continue to DB
	}

	var u URLRecord
	err := s.DB.GetContext(ctx, &u, "SELECT id, code, original_url, created_at, expires_at, clicks_count FROM urls WHERE code=$1", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	// set cache
	if s.Redis != nil {
		s.Redis.Set(ctx, cacheKey, u.Original, 24*time.Hour)
	}
	return &u, nil
}

func (s *Store) InsertClick(ctx context.Context, urlID int, code, ip, country, city, ua, referer string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO clicks (url_id, code, ip, country, city, user_agent, referer) VALUES ($1,$2,$3,$4,$5,$6,$7)`, urlID, code, ip, country, city, ua, referer)
	if err != nil {
		return err
	}
	// increment counter (atomic)
	_, err = s.DB.ExecContext(ctx, `UPDATE urls SET clicks_count = clicks_count + 1 WHERE id = $1`, urlID)
	return err
}
