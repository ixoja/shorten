package webserver

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ixoja/shorten/internal/webserver/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServer_Shorten(t *testing.T) {
	t.Run("400 Bad request error", func(t *testing.T) {
		ws := &Server{}
		req := httptest.NewRequest("POST", "http://example.com/foo", strings.NewReader(``))
		rr := httptest.NewRecorder()

		ws.Shorten(rr, req)
		res := rr.Result()
		assert.Equal(t, "400 Bad Request", res.Status)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("500 Internal Server Error", func(t *testing.T) {
		apiURL := "short.en/"
		client := mocks.HTTPClient{}
		ws := &Server{ApiURL: apiURL, Client: &client}
		client.On("Post", apiURL, url, "http://google.com").
			Return(nil, errors.New("some err"))
		req := httptest.NewRequest("POST",
			"http://example.com/foo",
			strings.NewReader(`url=http%3A%2F%2Fgoogle.com`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		rr := httptest.NewRecorder()

		ws.Shorten(rr, req)
		res := rr.Result()
		assert.Equal(t, "500 Internal Server Error", res.Status)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		client.AssertExpectations(t)
	})

	t.Run("302 Found success", func(t *testing.T) {
		apiURL := "short.en/"
		client := mocks.HTTPClient{}
		ws := &Server{ApiURL: apiURL, Client: &client}
		response := &http.Response{
			Status: "302 Found",
			StatusCode: http.StatusFound,
			Header:     make(map[string][]string),
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
		}
		client.On("Post", apiURL, url, "http://google.com").Return(response, nil)
		req := httptest.NewRequest("POST",
			"http://example.com/foo",
			strings.NewReader(`url=http%3A%2F%2Fgoogle.com`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		rr := httptest.NewRecorder()

		ws.Shorten(rr, req)
		res := rr.Result()
		assert.Equal(t, response.Status, res.Status)
		assert.Equal(t, response.StatusCode, res.StatusCode)
		client.AssertExpectations(t)
	})
}
