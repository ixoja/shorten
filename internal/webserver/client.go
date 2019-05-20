package webserver

import (
	"bytes"
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
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
