package model

import (
	"errors"
	"time"
)

type StoredURL struct {
	ID        string
	LongURL   string
	CreatedAt time.Time
}

var (
	ErrNotFound        = errors.New("not found")
	ErrStorageInternal = errors.New("storage internal")
	ErrEmptyArgument   = errors.New("argument cannot be empty")
)