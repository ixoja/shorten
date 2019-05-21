package webserver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Server struct {
	Client HTTPClient
	ApiURL string
}

//go:generate mockery -case=underscore -name HTTPClient
type HTTPClient interface {
	Post(url, key, value string) (*http.Response, error)
}

const url = "url"

func New(client HTTPClient, apiURL string) *Server {
	return &Server{Client: client, ApiURL:apiURL}
}

func (s *Server) Shorten(w http.ResponseWriter, r *http.Request) {
	val, err := extractValue(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.Client.Post(s.ApiURL, url, val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	fmt.Print(w, body)
}

func extractValue(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}

	vals, ok := r.Form[url]
	if !ok {
		return "", errors.New("no url in request")
	}
	val := vals[0]
	return val, nil
}
