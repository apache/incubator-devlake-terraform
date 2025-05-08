// Copyright (c) HashiCorp, Inc.

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// CreateApiKey - Creates new apikey.
func (c *Client) CreateApiKey(apiKeyCreate ApiKeyCreate) (*ApiKey, error) {
	rb, err := json.Marshal(apiKeyCreate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api-keys", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiKey := ApiKey{}
	err = json.Unmarshal(body, &apiKey)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

// ReadApiKeys - Returns list of apikeys.
func (c *Client) ReadApiKeys() ([]ApiKey, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api-keys", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiKeys := []ApiKey{}
	err = json.Unmarshal(body, &struct {
		ApiKeys *[]ApiKey `json:"apikeys"`
	}{
		ApiKeys: &apiKeys,
	})
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// DeleteApiKey - Deletes an apikey.
func (c *Client) DeleteApiKey(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api-keys/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	jsonBody := JsonBody{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return err
	}
	if !jsonBody.Success {
		return errors.New(jsonBody.Message)
	}

	return nil
}
