package webserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/ixoja/shorten/internal/grpcapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Server struct {
	Client grpcapi.ShortenServiceClient
	MyURL  string
}

//go:generate mockery -case=underscore -dir=../grpcapi -name ShortenServiceClient
const (
	url = "url"
)

func New(client grpcapi.ShortenServiceClient, myURL string) *Server {
	return &Server{Client: client, MyURL: myURL}
}

func (s *Server) Shorten(w http.ResponseWriter, r *http.Request) {
	longURL, err := extractValue(r, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.Client.Shorten(context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL})
	switch status.Code(err) {
	case codes.OK:
		shortURL := fmt.Sprintf("%s/to?%s", s.MyURL, resp.Hash)
		w.WriteHeader(http.StatusOK)
		fmt.Print(w, shortURL)
	case codes.InvalidArgument:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case codes.Internal:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var errNoURL = errors.New("no url in request")

func extractValue(r *http.Request, key string) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}

	vals, ok := r.Form[key]
	if !ok || vals[0] == "" {
		return "", errNoURL
	}
	val := vals[0]
	return val, nil
}

func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.RawQuery
	if hash == "" {
		http.Error(w, "no hash provided", http.StatusBadRequest)
		return
	}

	resp, err := s.Client.RedirectURL(context.Background(), &grpcapi.RedirectURLRequest{Hash: hash})
	switch status.Code(err) {
	case codes.OK:
		http.Redirect(w, r, resp.LongUrl, http.StatusFound)
	case codes.InvalidArgument:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case codes.Internal:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
