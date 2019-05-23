package model

import "time"

type StoredURL struct {
	ID        string
	LongURL   string
	CreatedAt time.Time
}

