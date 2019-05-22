package handler

import (
	"context"
	"github.com/ixoja/shorten/internal/grpcapi"
)

type Service struct {
	controller Controller
}

type Controller interface {
	Shorten(url string) (string, error)
	RedirectURL(hash string) (string, error)
}

func (s *Service) Shorten(ctx context.Context, in *grpcapi.ShortenRequest) (*grpcapi.ShortenResponse, error) {
	return nil, nil
}

func (s *Service) RedirectURL(ctx context.Context, in *grpcapi.RedirectURLRequest) (*grpcapi.RedirectURLResponse, error) {
	return nil, nil
}
