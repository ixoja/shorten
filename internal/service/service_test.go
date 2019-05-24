package service

import (
	"context"
	"github.com/ixoja/shorten/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/icrowley/fake"
	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Shorten(t *testing.T) {
	for name, tc := range map[string]struct {
		have error
		want codes.Code
	}{
		"error internal": {
			have: model.ErrStorageInternal,
			want: codes.Internal,
		},
		"error empty argument": {
			have: model.ErrEmptyArgument,
			want: codes.InvalidArgument,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := mocks.Controller{}
			s := Service{Controller: &c}
			longURL := fake.DomainName()

			c.On("Shorten", longURL).Return("", tc.have)
			_, err := s.Shorten(context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL})
			assertCodes(t, tc.want, err)
			c.AssertExpectations(t)
		})
	}

	t.Run("success", func(t *testing.T) {
		c := mocks.Controller{}
		s := Service{Controller: &c}
		longURL := fake.DomainName()
		hash := fake.CharactersN(5)

		c.On("Shorten", longURL).Return(hash, nil)
		res, err := s.Shorten(context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL})
		require.NoError(t, err)
		assert.Equal(t, &grpcapi.ShortenResponse{Hash: hash}, res)
		c.AssertExpectations(t)
	})

}

func TestService_RedirectURL(t *testing.T) {
	for name, tc := range map[string]struct {
		have error
		want codes.Code
	}{
		"error not found": {
			have: model.ErrNotFound,
			want: codes.NotFound,
		},
		"error internal": {
			have: model.ErrStorageInternal,
			want: codes.Internal,
		},
		"error empty argument": {
			have: model.ErrEmptyArgument,
			want: codes.InvalidArgument,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := mocks.Controller{}
			s := Service{Controller: &c}
			hash := fake.CharactersN(5)

			c.On("RedirectURL", hash).Return("", tc.have)
			_, err := s.RedirectURL(context.Background(), &grpcapi.RedirectURLRequest{Hash: hash})
			assertCodes(t, tc.want, err)
			c.AssertExpectations(t)
		})
	}

	t.Run("success", func(t *testing.T) {
		c := mocks.Controller{}
		s := Service{Controller: &c}
		hash := fake.CharactersN(5)
		longURL := fake.DomainName()

		c.On("RedirectURL", hash).Return(longURL, nil)
		res, err := s.RedirectURL(context.Background(), &grpcapi.RedirectURLRequest{Hash: hash})
		require.NoError(t, err)
		assert.Equal(t, &grpcapi.RedirectURLResponse{LongUrl: longURL}, res)
		c.AssertExpectations(t)
	})
}

func assertCodes(t *testing.T, code codes.Code, err error) {
	s, _ := status.FromError(err)
	assert.Equal(t, code, s.Code())
}
