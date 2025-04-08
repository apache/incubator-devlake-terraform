// Copyright (c) HashiCorp, Inc.

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetBitbucketServerConnection - Returns bitbucket server connection.
func (c *Client) GetBitbucketServerConnection(id string) (*BitbucketServerConnection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	connection := BitbucketServerConnection{}
	err = json.Unmarshal(body, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

// CreateBitbucketServerConnection - Creates new bitbucket server connection.
func (c *Client) CreateBitbucketServerConnection(connection BitbucketServerConnection) (*BitbucketServerConnection, error) {
	rb, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/plugins/bitbucket_server/connections", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	createdConnection := BitbucketServerConnection{}
	err = json.Unmarshal(body, &createdConnection)
	if err != nil {
		return nil, err
	}

	return &createdConnection, nil
}

// UpdateBitbucketServerConnection - Updates bitbucket server connection.
func (c *Client) UpdateBitbucketServerConnection(id string, connection BitbucketServerConnection) (*BitbucketServerConnection, error) {
	rb, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s", c.HostURL, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedConnection := BitbucketServerConnection{}
	err = json.Unmarshal(body, &updatedConnection)
	if err != nil {
		return nil, err
	}

	return &updatedConnection, nil
}

// DeleteBitbucketServerConnection - Deletes an apikey.
func (c *Client) DeleteBitbucketServerConnection(id string) error {
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
