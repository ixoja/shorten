package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ixoja/shorten/internal/restapi/operations"
	"github.com/mikkeloscar/gin-swagger/api"
)

type Service struct {
}

func (s *Service) Healthy() bool {
	return true
}
func (s *Service) Redirect(ctx *gin.Context, params *operations.RedirectParams) *api.Response {
	return &api.Response{}
}
func (s *Service) Shorten(ctx *gin.Context, params *operations.ShortenParams) *api.Response {
	return &api.Response{}
}
