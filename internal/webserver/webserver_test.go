package webserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/ixoja/shorten/internal/grpcapi"
	"github.com/ixoja/shorten/internal/webserver/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const httpPrefix = "http://"

func TestServer_Shorten(t *testing.T) {
	t.Run("400 Bad request error", func(t *testing.T) {
		t.Run("http bad request", func(t *testing.T) {
			ws := &Server{}
			req := httptest.NewRequest("POST", httpPrefix+fake.DomainName(), strings.NewReader(``))
			rr := httptest.NewRecorder()

			ws.Shorten(rr, req)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})

		t.Run("grpc bad request", func(t *testing.T) {
			client := mocks.ShortenServiceClient{}
			ws := &Server{Client: &client}
			longURL := fake.DomainName()
			client.On("Shorten", context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL}).
				Return(nil, status.Error(codes.InvalidArgument, "bad arg"))
			req := preparePostRequest(longURL)
			rr := httptest.NewRecorder()

			ws.Shorten(rr, req)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			client.AssertExpectations(t)
		})
	})

	t.Run("500 Internal Server Error", func(t *testing.T) {
		client := mocks.ShortenServiceClient{}
		ws := &Server{Client: &client}
		longURL := fake.DomainName()
		client.On("Shorten", context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL}).
			Return(nil, status.Error(codes.Internal, "internal"))
		req := preparePostRequest(longURL)
		rr := httptest.NewRecorder()

		ws.Shorten(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		client.AssertExpectations(t)
	})

	t.Run("200 OK success", func(t *testing.T) {
		client := mocks.ShortenServiceClient{}
		myURL := fake.DomainName()
		ws := &Server{MyURL: myURL, Client: &client}
		longURL := fake.DomainName()
		hash := fake.CharactersN(5)
		client.On("Shorten", context.Background(), &grpcapi.ShortenRequest{LongUrl: longURL}).
			Return(&grpcapi.ShortenResponse{Hash: hash}, nil)
		req := preparePostRequest(longURL)
		rr := httptest.NewRecorder()

		ws.Shorten(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, shortURL(myURL, hash), rr.Body.String())
		client.AssertExpectations(t)
	})
}

func preparePostRequest(longURL string) *http.Request {
	req := httptest.NewRequest("POST",
		httpPrefix+fake.DomainName(), strings.NewReader(fmt.Sprintf(`url=%s`, longURL)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	return req
}

func Test_extractValue(t *testing.T) {
	domain := fake.DomainName()
	for name, tc := range map[string]struct {
		req *http.Request
		res string
		err error
	}{
		"success": {
			req: httptest.NewRequest("POST",
				httpPrefix+fake.DomainName(), strings.NewReader(fmt.Sprintf(`url=%s`, domain))),
			res: domain,
			err: nil,
		},
		"empty url error": {
			req: httptest.NewRequest("POST",
				httpPrefix+fake.DomainName(), strings.NewReader(`url=`)),
			res: "",
			err: errNoURL,
		},
		"no url parameter error": {
			req: httptest.NewRequest("POST",
				httpPrefix+fake.DomainName(), strings.NewReader(``)),
			res: "",
			err: errNoURL,
		},
	} {
		t.Run(name, func(t *testing.T) {
			tc.req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			res, err := extractValue(tc.req, urlConst)
			assert.Equal(t, tc.res, res)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestServer_Redirect(t *testing.T) {
	t.Run("400 Bad request error", func(t *testing.T) {
		t.Run("http bad request", func(t *testing.T) {
			ws := &Server{}
			req := &http.Request{URL: &url.URL{}}
			rr := httptest.NewRecorder()

			ws.Redirect(rr, req)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})

		t.Run("grpc bad request", func(t *testing.T) {
			client := mocks.ShortenServiceClient{}
			ws := &Server{Client: &client}
			hash := fake.CharactersN(5)
			client.On("RedirectURL", context.Background(), &grpcapi.RedirectURLRequest{Hash: hash}).
				Return(nil, status.Error(codes.InvalidArgument, "bad arg"))
			req := &http.Request{URL: &url.URL{RawQuery: hash}}
			rr := httptest.NewRecorder()

			ws.Redirect(rr, req)
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			client.AssertExpectations(t)
		})
	})

	t.Run("grpc error", func(t *testing.T) {
		for name, tc := range map[string]struct {
			grpcCode codes.Code
			httpCode int
		}{
			"500 internal": {
				grpcCode: codes.Internal,
				httpCode: http.StatusInternalServerError,
			},
			"400 invalid argument": {
				grpcCode: codes.InvalidArgument,
				httpCode: http.StatusBadRequest,
			},
			"404 not found": {
				grpcCode: codes.NotFound,
				httpCode: http.StatusNotFound,
			},
		} {
			t.Run(name, func(t *testing.T) {
				client := mocks.ShortenServiceClient{}
				ws := &Server{Client: &client}
				hash := fake.CharactersN(5)
				client.On("RedirectURL", context.Background(), &grpcapi.RedirectURLRequest{Hash: hash}).
					Return(nil, status.Error(tc.grpcCode, "internal"))
				req := &http.Request{URL: &url.URL{RawQuery: hash}}
				rr := httptest.NewRecorder()

				ws.Redirect(rr, req)
				assert.Equal(t, tc.httpCode, rr.Code)
				client.AssertExpectations(t)
			})
		}
	})

	t.Run("200 OK success", func(t *testing.T) {
		client := mocks.ShortenServiceClient{}
		ws := &Server{Client: &client}
		longURL := httpPrefix + fake.DomainName()
		hash := fake.CharactersN(5)
		client.On("RedirectURL", context.Background(), &grpcapi.RedirectURLRequest{Hash: hash}).
			Return(&grpcapi.RedirectURLResponse{LongUrl: longURL}, nil)
		req := &http.Request{URL: &url.URL{RawQuery: hash}}
		rr := httptest.NewRecorder()

		ws.Redirect(rr, req)
		assert.Equal(t, http.StatusFound, rr.Code)
		assert.Equal(t, longURL, rr.Header().Get("Location"))
		client.AssertExpectations(t)
	})
}
