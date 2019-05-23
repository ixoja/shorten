package controller

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type Controller struct {
	cache Storage
	storage Storage
}

type StoredURL struct {
	ID string
	LongURL string
	CreatedAt time.Time
}

type Storage interface {
	Save(stored *StoredURL) (*StoredURL, error)
	Delete(id string) error
	Get(id string) (*StoredURL, bool, error)
	GetByURL(longURL string) (*StoredURL, bool, error)
	EvictOlder(timestamp time.Time) error
}

func (c *Controller) Shorten(longURL string) (string, error) {
	longURL = toFullURL(longURL)
	if stored, ok, err := c.lookupByURL(longURL); err != nil {
		return "", status.Error(codes.Internal, err.Error())
	} else if ok {
		return stored.ID, nil
	}

	stored, err := c.save(longURL)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return stored.ID, nil
}

func (c *Controller) lookupByURL(longURL string) (*StoredURL, bool, error) {
	if stored, ok, err := c.cache.GetByURL(longURL); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by url from cache")
	} else if ok {
		return stored, ok, nil
	}

	if stored, ok, err := c.cache.GetByURL(longURL); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by url from storage")
	} else if ok {
		return stored, ok, nil
	}

	return nil, false, nil
}

func (c *Controller) save(longURL string) (*StoredURL, error) {
	stored, err := c.storage.Save(&StoredURL{LongURL: longURL})
	if err != nil {
		return nil, errors.Wrap(err, "failed to save into storage")
	}

	_, err = c.cache.Save(stored)
	if err != nil {
		if err := c.storage.Delete(stored.ID); err != nil {
			return nil, errors.Wrap(err, "failed to delete from storage")
		}
		return nil, errors.Wrap(err, "failed to save into cache")
	}

	return stored, nil
}

func (c *Controller) RedirectURL(hash string) (string, error) {
	stored, ok, err := c.lookupByID(hash)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}
	if ok {
		return stored.ID, nil
	}

	return "", status.Error(codes.NotFound, "has not found")
}

func (c *Controller) lookupByID(id string) (*StoredURL, bool, error) {
	if stored, ok, err := c.cache.Get(id); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by id from cache")
	} else if ok {
		return stored, ok, nil
	}

	if stored, ok, err := c.cache.Get(id); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by id from storage")
	} else if ok {
		return stored, ok, nil
	}

	return nil, false, nil
}

const (
	http = "http://"
	https = "https://"
	ftp = "ftp://"
)

func toFullURL(s string) string {
	if !strings.HasPrefix(s, http) ||
		!strings.HasPrefix(s, https) ||
		!strings.HasPrefix(s, ftp) {
		return http + s
	}
	return s
}