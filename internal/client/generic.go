// Copyright (c) HashiCorp, Inc.

package client

import (
	"encoding/json"
	"net/http"
	"strings"
)

// create - Generic wrapper for POST requests
func create[T any](c *Client, url string, reqObj T) (*T, error) {
	rb, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var respObj T
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return nil, err
	}

	return &respObj, nil
}

// read - Generic wrapper for GET requests
func read[T any](c *Client, url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var respObj T
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return nil, err
	}

	return &respObj, nil
}

// update - Generic wrapper for PATCH requests
func update[T any](c *Client, url string, reqObj T) (*T, error) {
	rb, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var respObj T
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return nil, err
	}

	return &respObj, nil
}

// delete - Generic wrapper for DELETE requests
func delete(c *Client, url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
