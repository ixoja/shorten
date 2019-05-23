package service

import (
	"context"
	"github.com/pkg/errors"
	"testing"

	"github.com/icrowley/fake"
	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Shorten(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		c := mocks.Controller{}
		s := Service{controller: &c}
		retErr := errors.New("some error")
		longURL := fake.DomainName()

		c.On("Shorten", longURL).Return("", retErr)
		_, err := s.Shorten(context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL})
		assert.Equal(t, retErr, errors.Cause(err))
		c.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		c := mocks.Controller{}
		s := Service{controller: &c}
		longURL := fake.DomainName()
		shortURL := fake.CharactersN(5)

		c.On("Shorten", longURL).Return(shortURL, nil)
		res, err := s.Shorten(context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL})
		require.NoError(t, err)
		assert.Equal(t, &grpcapi.ShortenResponse{ShortUrl: shortURL}, res)
		c.AssertExpectations(t)
	})

}

func TestService_RedirectURL(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		c := mocks.Controller{}
		s := Service{controller: &c}
		hash := fake.CharactersN(5)
		retErr := errors.New("some error")

		c.On("RedirectURL", hash).Return("", retErr)
		_, err := s.RedirectURL(context.Background(), &grpcapi.RedirectURLRequest{Hash:hash})
		assert.Equal(t, retErr, errors.Cause(err))
		c.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		c := mocks.Controller{}
		s := Service{controller: &c}
		hash := fake.CharactersN(5)
		longURL := fake.DomainName()

		c.On("RedirectURL", hash).Return(longURL, nil)
		res, err := s.RedirectURL(context.Background(), &grpcapi.RedirectURLRequest{Hash:hash})
		require.NoError(t, err)
		assert.Equal(t, &grpcapi.RedirectURLResponse{LongUrl: longURL}, res)
		c.AssertExpectations(t)
	})
}
