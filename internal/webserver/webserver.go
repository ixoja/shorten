package webserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ixoja/shorten/internal/grpcapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Client grpcapi.ShortenServiceClient
	MyURL  string
}

//go:generate mockery -case=underscore -dir=../grpcapi -name ShortenServiceClient
const (
	urlConst = "url"
)

func New(client grpcapi.ShortenServiceClient, myURL string) *Server {
	return &Server{Client: client, MyURL: myURL}
}

func (s *Server) Shorten(w http.ResponseWriter, r *http.Request) {
	longURL, err := extractValue(r, urlConst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.Client.Shorten(context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL})
	switch status.Code(err) {
	case codes.OK:
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, shortURL(s.MyURL, resp.Hash)); err != nil {
			log.Println(err.Error())
		}
	case codes.InvalidArgument:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case codes.Internal:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func shortURL(myURL, hash string) string {
	return myURL + "/to?" + hash
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
	case codes.NotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case codes.Internal:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
