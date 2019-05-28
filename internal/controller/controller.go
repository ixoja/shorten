package controller

import (
	"github.com/ixoja/shorten/internal/model"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Controller struct {
	Cache   Storage
	Storage Storage
}

func New(cache, storage Storage) *Controller {
	return &Controller{Cache: cache, Storage: storage}
}

//go:generate mockery -case=underscore -name Storage
type Storage interface {
	Save(stored *model.StoredURL) (*model.StoredURL, error)
	Delete(id string) error
	Get(id string) (*model.StoredURL, bool, error)
	GetByURL(longURL string) (*model.StoredURL, bool, error)
	EvictOlder(timestamp time.Time) error
}

func (c *Controller) Shorten(longURL string) (string, error) {
	if longURL == "" {
		return "", errors.Wrap(model.ErrEmptyArgument, "long url")
	}

	longURL = toFullURL(longURL)
	if stored, ok, err := c.lookupByURL(longURL); err != nil {
		return "", errors.Wrap(model.ErrStorageInternal, err.Error())
	} else if ok {
		return stored.ID, nil
	}

	stored, err := c.save(longURL)
	if err != nil {
		return "", errors.Wrap(model.ErrStorageInternal, err.Error())
	}

	return stored.ID, nil
}

func (c *Controller) lookupByURL(longURL string) (*model.StoredURL, bool, error) {
	if stored, ok, err := c.Cache.GetByURL(longURL); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by url from cache")
	} else if ok {
		return stored, ok, nil
	}

	if stored, ok, err := c.Storage.GetByURL(longURL); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by url from storage")
	} else if ok {
		return stored, ok, nil
	}

	return nil, false, nil
}

func (c *Controller) save(longURL string) (*model.StoredURL, error) {
	stored, err := c.Storage.Save(&model.StoredURL{LongURL: longURL})
	if err != nil {
		return nil, errors.Wrap(err, "failed to save into storage")
	}

	_, err = c.Cache.Save(stored)
	if err != nil {
		if err := c.Storage.Delete(stored.ID); err != nil {
			return nil, errors.Wrap(err, "failed to delete from storage")
		}
		return nil, errors.Wrap(err, "failed to save into cache")
	}

	return stored, nil
}

func (c *Controller) RedirectURL(hash string) (string, error) {
	if hash == "" {
		return "", errors.Wrap(model.ErrEmptyArgument, "hash")
	}

	stored, ok, err := c.lookupByID(hash)
	if err != nil {
		return "", errors.Wrap(model.ErrStorageInternal, err.Error())
	}
	if ok {
		return stored.LongURL, nil
	}

	return "", errors.Wrap(model.ErrNotFound, "hashZ")
}

func (c *Controller) lookupByID(id string) (*model.StoredURL, bool, error) {
	if stored, ok, err := c.Cache.Get(id); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by id from cache")
	} else if ok {
		return stored, ok, nil
	}

	if stored, ok, err := c.Storage.Get(id); err != nil {
		return nil, false, errors.Wrap(err, "failed to get by id from storage")
	} else if ok {
		return stored, ok, nil
	}

	return nil, false, nil
}

const (
	http  = "http://"
	https = "https://"
	ftp   = "ftp://"
)

func toFullURL(s string) string {
	if !strings.HasPrefix(s, http) &&
		!strings.HasPrefix(s, https) &&
		!strings.HasPrefix(s, ftp) {
		return http + s
	}
	return s
}
