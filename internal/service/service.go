package service

import (
	"context"
	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Controller Controller
}

func New(controller Controller) *Service {
	return &Service{Controller: controller}
}

//go:generate mockery -case=underscore -name Controller
type Controller interface {
	Shorten(url string) (string, error)
	RedirectURL(hash string) (string, error)
}

func (s *Service) Shorten(ctx context.Context, r *grpcapi.ShortenRequest) (*grpcapi.ShortenResponse, error) {
	switch hash, err := s.Controller.Shorten(r.LongUrl); errors.Cause(err) {
	case nil:
		return &grpcapi.ShortenResponse{Hash: hash}, nil
	case model.ErrEmptyArgument:
		return &grpcapi.ShortenResponse{}, status.Error(codes.InvalidArgument, err.Error())
	default:
		return &grpcapi.ShortenResponse{}, status.Error(codes.Internal, err.Error())
	}
}

func (s *Service) RedirectURL(ctx context.Context, r *grpcapi.RedirectURLRequest) (*grpcapi.RedirectURLResponse, error) {
	switch url, err := s.Controller.RedirectURL(r.Hash); errors.Cause(err) {
	case nil:
		return &grpcapi.RedirectURLResponse{LongUrl: url}, nil
	case model.ErrEmptyArgument:
		return &grpcapi.RedirectURLResponse{}, status.Error(codes.InvalidArgument, err.Error())
	case model.ErrNotFound:
		return &grpcapi.RedirectURLResponse{}, status.Error(codes.NotFound, err.Error())
	default:
		return &grpcapi.RedirectURLResponse{}, status.Error(codes.Internal, err.Error())
	}
}
