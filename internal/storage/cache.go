package storage

import (
	"github.com/ixoja/shorten/internal/model"
	"sync"
	"time"
)

type Cache struct {
	urls struct {
		sync.RWMutex
		m map[string]*model.StoredURL
	}
}

func NewCache() *Cache {
	m := make(map[string]*model.StoredURL)
	return &Cache{urls: struct {
		sync.RWMutex
		m map[string]*model.StoredURL
	}{m: m}}
}

func (c *Cache) Save(stored *model.StoredURL) (*model.StoredURL, error) {
	c.urls.Lock()
	c.urls.m[stored.ID] = stored
	c.urls.Unlock()
	return stored, nil
}
func (c *Cache) Delete(key string) error {
	c.urls.Lock()
	delete(c.urls.m, key)
	c.urls.Unlock()
	return nil
}

func (c *Cache) Get(key string) (*model.StoredURL, bool, error) {
	c.urls.RLock()
	stored, ok := c.urls.m[key]
	c.urls.RUnlock()
	if ok {
		return stored, true, nil
	}
	return nil, false, nil
}

func (c *Cache) GetByURL(longURL string) (*model.StoredURL, bool, error) {
	c.urls.RLock()
	for _, stored := range c.urls.m {
		if stored.LongURL == longURL {
			return stored, true, nil
		}
	}
	c.urls.RUnlock()
	return nil, false, nil
}

func (c *Cache) EvictOlder(timestamp time.Time) error {
	return nil
}
