// Copyright (c) HashiCorp, Inc.

package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Devlake URL.
const HostURL string = "http://localhost:8080"

// Client - Holds backend data.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient - Create new client.
func NewClient(host, token *string) (*Client, error) {
	// If token not provided, return error
	if token == nil {
		return nil, errors.New("no api token provided")
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Devlake URL
		HostURL: HostURL,
		Token:   *token,
	}

	if host != nil {
		c.HostURL = *host
	}

	return &c, nil
}

// doRequest - Query the devlake backend.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
