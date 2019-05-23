package controller

import (
	"github.com/ixoja/shorten/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/ixoja/shorten/internal/controller/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestController_Shorten(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Run("lookup in cache", func(t *testing.T) {
			cache := mocks.Storage{}
			c := Controller{Cache: &cache}
			longURL := fake.DomainName()
			retErr := errors.New("cache err")
			cache.On("GetByURL", toFullURL(longURL)).Return(nil, false, retErr)

			_, err := c.Shorten(longURL)
			code, _ := status.FromError(err)
			assert.Equal(t, codes.Internal, code.Code())
			cache.AssertExpectations(t)
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

		t.Run("found in storage", func(t *testing.T) {
			cache := mocks.Storage{}
			storage := mocks.Storage{}
			c := Controller{Cache: &cache, Storage: &storage}
			longURL := fake.DomainName()
			hash := fake.CharactersN(5)
			cache.On("GetByURL", toFullURL(longURL)).Return(nil, false, nil)
			storage.On("GetByURL", toFullURL(longURL)).Return(&model.StoredURL{ID: hash}, true, nil)

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
			stored := &model.StoredURL{LongURL: fullURL, ID: hash, CreatedAt: time.Now()}
			storage.On("Save",&model.StoredURL{LongURL: fullURL}).Return(stored, nil)
			cache.On("Save",stored).Return(stored, nil)

			res, err := c.Shorten(longURL)
			require.NoError(t, err)
			assert.Equal(t, hash, res)
			mock.AssertExpectationsForObjects(t, &cache, &storage)
		})
	})
}
