package webserver

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
)

type Client struct {
	Client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{Client: client}
}

func (c Client) Post(url, key, value string) (*http.Response, error) {
	var jsonStr = []byte(`{"` + key + `":"` + value + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, errors.Wrap(err, "error when creating post request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error when calling post request")
	}

	if err := resp.Body.Close(); err != nil {
		return nil, errors.Wrap(err, "error when closing post request")
	}

	return resp, nil
}
