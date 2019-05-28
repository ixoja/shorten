package webserver

import (
	"bytes"
	"fmt"
	"github.com/icrowley/fake"
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
		myURL := fake.DomainName()
		client := mocks.ShortenServiceClient{}
		ws := &Server{MyURL: myURL, Client: &client}

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
		myURL := "short.en/"
		client := mocks.ShortenServiceClient{}
		ws := &Server{MyURL: myURL, Client: &client}
		response := &http.Response{
			Status:     "302 Found",
			StatusCode: http.StatusFound,
			Header:     make(map[string][]string),
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
		}
		client.On("Post", myURL, url, "http://google.com").Return(response, nil)
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

func Test_extractValue(t *testing.T) {
	url := fake.DomainName()
	for name, tc := range map[string]struct {
		req *http.Request
		res string
		err error
	}{
		"success": {
			req: httptest.NewRequest("POST",
				"http://example.com/foo", strings.NewReader(fmt.Sprintf(`url=%s`, url))),
			res: url,
			err: nil,
		},
		"empty url error": {
			req: httptest.NewRequest("POST",
				"http://example.com/foo", strings.NewReader(`url=`)),
			res: "",
			err: errNoURL,
		},
		"no url parameter error": {
			req: httptest.NewRequest("POST",
				"http://example.com/foo", strings.NewReader(``)),
			res: "",
			err: errNoURL,
		},
	}{
		t.Run(name, func(t *testing.T) {
			tc.req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			res, err := extractValue(tc.req, url)
			assert.Equal(t, tc.res, res)
			assert.Equal(t, tc.err, err)
		})
	}
}
