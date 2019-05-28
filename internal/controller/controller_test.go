package controller

import (
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/ixoja/shorten/internal/controller/mocks"
	"github.com/ixoja/shorten/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestController_Shorten(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Run("empty url", func(t *testing.T) {
			c := Controller{}
			_, err := c.Shorten("")
			assert.Equal(t, model.ErrEmptyArgument, errors.Cause(err))
		})

		t.Run("get from cache", func(t *testing.T) {
			cache := mocks.Storage{}
			c := Controller{Cache: &cache}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			retErr := errors.New("cache err")
			cache.On("GetByURL", fullURL).Return(nil, false, retErr)

			_, err := c.Shorten(longURL)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			cache.AssertExpectations(t)
		})

		t.Run("get from storage", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			retErr := errors.New("storage err")
			cache.On("GetByURL", fullURL).Return(nil, false, nil)
			storage.On("GetByURL", fullURL).Return(nil, false, retErr)

			_, err := c.Shorten(longURL)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("save to storage", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			cache.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			storage.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			retErr := errors.New("storage err")
			storage.On("Save", &model.StoredURL{LongURL: fullURL}).Return(nil, retErr)

			_, err := c.Shorten(longURL)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("save to cache successful delete", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			cache.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			storage.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			hash := fake.CharactersN(5)
			stored := &model.StoredURL{LongURL: fullURL, ID: hash, LastAccess: time.Now()}
			storage.On("Save", &model.StoredURL{LongURL: fullURL}).Return(stored, nil)
			retErr := errors.New("cache err")
			cache.On("Save", stored).Return(nil, retErr)
			storage.On("Delete", hash).Return(nil)

			_, err := c.Shorten(longURL)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("save to cache delete failed", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			cache.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			storage.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			hash := fake.CharactersN(5)
			stored := &model.StoredURL{LongURL: fullURL, ID: hash, LastAccess: time.Now()}
			storage.On("Save", &model.StoredURL{LongURL: fullURL}).Return(stored, nil)
			cacheErr := errors.New("cache err")
			cache.On("Save", stored).Return(nil, cacheErr)
			storageErr := errors.New("storage err")
			storage.On("Delete", hash).Return(storageErr)

			_, err := c.Shorten(longURL)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})
	})

	t.Run("success", func(t *testing.T) {
		t.Run("found in cache", func(t *testing.T) {
			cache := mocks.Storage{}
			c := Controller{Cache: &cache}
			longURL := fake.DomainName()
			hash := fake.CharactersN(5)
			cache.On("GetByURL", toFullURL(longURL)).Return(&model.StoredURL{ID: hash}, true, nil)

			res, err := c.Shorten(longURL)
			require.NoError(t, err)
			assert.Equal(t, hash, res)
			cache.AssertExpectations(t)
		})

		t.Run("found in storage saved to cache", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			hash := fake.CharactersN(5)
			stored := &model.StoredURL{ID: hash}
			cache.On("GetByURL", fullURL).Return(nil, false, nil)
			storage.On("GetByURL", fullURL).Return(stored, true, nil)
			cache.On("Save", stored).Return(stored, nil)

			res, err := c.Shorten(longURL)
			require.NoError(t, err)
			assert.Equal(t, hash, res)
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("found in storage failed to saved to cache", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			hash := fake.CharactersN(5)
			stored := &model.StoredURL{ID: hash}
			cache.On("GetByURL", fullURL).Return(nil, false, nil)
			storage.On("GetByURL", fullURL).Return(stored, true, nil)
			cache.On("Save", stored).Return(nil, errors.New("cache error"))

			res, err := c.Shorten(longURL)
			require.NoError(t, err)
			assert.Equal(t, hash, res)
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("not found in storage, created successfully", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			hash := fake.CharactersN(5)
			cache.On("GetByURL", fullURL).Return(nil, false, nil)
			storage.On("GetByURL", fullURL).Return(nil, false, nil)
			stored := &model.StoredURL{LongURL: fullURL, ID: hash, LastAccess: time.Now()}
			storage.On("Save", &model.StoredURL{LongURL: fullURL}).Return(stored, nil)
			cache.On("Save", stored).Return(stored, nil)

			res, err := c.Shorten(longURL)
			require.NoError(t, err)
			assert.Equal(t, hash, res)
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})
	})
}

func TestController_RedirectURL(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Run("empty hash", func(t *testing.T) {
			c := Controller{}
			_, err := c.RedirectURL("")
			assert.Equal(t, model.ErrEmptyArgument, errors.Cause(err))
		})

		t.Run("get from cache", func(t *testing.T) {
			cache := mocks.Storage{}
			c := Controller{Cache: &cache}
			hash := fake.CharactersN(5)
			retErr := errors.New("cache err")
			cache.On("Get", hash).Return(nil, false, retErr)

			_, err := c.RedirectURL(hash)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			cache.AssertExpectations(t)
		})

		t.Run("get from storage", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			hash := fake.CharactersN(5)
			retErr := errors.New("storage err")
			cache.On("Get", hash).Return(nil, false, nil)
			storage.On("Get", hash).Return(nil, false, retErr)

			_, err := c.RedirectURL(hash)
			assert.Equal(t, model.ErrStorageInternal, errors.Cause(err))
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("not found", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			hash := fake.CharactersN(5)
			cache.On("Get", hash).Return(nil, false, nil)
			storage.On("Get", hash).Return(nil, false, nil)

			_, err := c.RedirectURL(hash)
			assert.Equal(t, model.ErrNotFound, errors.Cause(err))
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})
	})

	t.Run("success", func(t *testing.T) {
		t.Run("found in cache", func(t *testing.T) {
			cache := mocks.Storage{}
			c := Controller{Cache: &cache}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			hash := fake.CharactersN(5)
			cache.On("Get", hash).Return(&model.StoredURL{ID: hash, LongURL: fullURL}, true, nil)

			res, err := c.RedirectURL(hash)
			require.NoError(t, err)
			assert.Equal(t, fullURL, res)
			cache.AssertExpectations(t)
		})

		t.Run("found in storage saved to cache", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			hash := fake.CharactersN(5)
			stored := &model.StoredURL{ID: hash, LongURL: fullURL}
			cache.On("Get", hash).Return(nil, false, nil)
			storage.On("Get", hash).Return(stored, true, nil)
			cache.On("Save", stored).Return(stored, nil)

			res, err := c.RedirectURL(hash)
			require.NoError(t, err)
			assert.Equal(t, fullURL, res)
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})

		t.Run("found in storage failed to save to cache", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			fullURL := toFullURL(longURL)
			hash := fake.CharactersN(5)
			stored := &model.StoredURL{ID: hash, LongURL: fullURL}
			cache.On("Get", hash).Return(nil, false, nil)
			storage.On("Get", hash).Return(stored, true, nil)
			cache.On("Save", stored).Return(nil, errors.New("cache error"))

			res, err := c.RedirectURL(hash)
			require.NoError(t, err)
			assert.Equal(t, fullURL, res)
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})
	})
}

func Test_toFullURL(t *testing.T) {
	url := fake.DomainName()
	for name, tc := range map[string]struct {
		have string
		want string
	}{
		"http": {
			have: http + url,
			want: http + url,
		},
		"https": {
			have: https + url,
			want: https + url,
		},
		"ftp": {
			have: ftp + url,
			want: ftp + url,
		},
		"no prefix": {
			have: url,
			want: http + url,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, toFullURL(tc.have))
		})
	}
}
