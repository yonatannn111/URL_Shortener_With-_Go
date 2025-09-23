package storage

import (
    "github.com/jmoiron/sqlx"
)

type Store struct {
    DB *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
    return &Store{DB: db}
}

func (s *Store) SaveURL(code string, longURL string) error {
    _, err := s.DB.Exec("INSERT INTO urls (code, original_url) VALUES ($1, $2)", code, longURL)
    return err
}

func (s *Store) GetURL(code string) (string, error) {
    var longURL string
    err := s.DB.Get(&longURL, "SELECT original_url FROM urls WHERE code=$1", code)
    return longURL, err
}
