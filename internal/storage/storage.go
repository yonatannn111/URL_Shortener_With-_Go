package storage

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

func NewStore(db *sqlx.DB, rdb *redis.Client) *Store {
	return &Store{DB: db, RDB: rdb}
}

// Example methods — you can expand later
func (s *Store) SaveURL(ctx context.Context, code, url string) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO urls (code, long_url) VALUES ($1, $2)", code, url)
	return err
}

func (s *Store) GetURL(ctx context.Context, code string) (string, error) {
	var longURL string
	err := s.DB.GetContext(ctx, &longURL, "SELECT long_url FROM urls WHERE code=$1", code)
	return longURL, err
}
