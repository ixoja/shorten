package service

import (
	"context"
	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/pkg/errors"
)

type Service struct {
	controller Controller
}

//go:generate mockery -case=underscore -name Controller
type Controller interface {
	Shorten(url string) (string, error)
	RedirectURL(hash string) (string, error)
}

func (s *Service) Shorten(ctx context.Context, r *grpcapi.ShortenRequest) (*grpcapi.ShortenResponse, error) {
	shortURL, err := s.controller.Shorten(r.LongUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to shorten url")
	}

	return &grpcapi.ShortenResponse{ShortUrl: shortURL}, nil
}

func (s *Service) RedirectURL(ctx context.Context, r *grpcapi.RedirectURLRequest) (*grpcapi.RedirectURLResponse, error) {
	url, err := s.controller.RedirectURL(r.Hash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get long url")
	}

	return &grpcapi.RedirectURLResponse{LongUrl: url}, nil
}
