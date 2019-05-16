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
func (s *Service) GetHash(ctx *gin.Context, params *operations.GetHashParams) *api.Response {
	return &api.Response{}
}
func (s *Service) PostShorten(ctx *gin.Context, params *operations.PostShortenParams) *api.Response {
	return &api.Response{}
}