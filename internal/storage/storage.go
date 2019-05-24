package storage

import (
	"database/sql"
	"time"

	"github.com/ixoja/shorten/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Database sql.DB
}

func New(db sql.DB) *SQLite {
	return &SQLite{Database:db}
}

func (s *SQLite) Save(stored *model.StoredURL) (*model.StoredURL, error) {

	return nil, nil
}
func (s *SQLite) Delete(id string) error {
	return nil
}
func (s *SQLite) Get(id string) (*model.StoredURL, bool, error) {
	return nil, false, nil
}
func (s *SQLite) GetByURL(longURL string) (*model.StoredURL, bool, error) {
	return nil, false, nil
}
func (s *SQLite) EvictOlder(timestamp time.Time) error {
	return nil
}
