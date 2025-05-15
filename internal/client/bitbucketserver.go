// Copyright (c) HashiCorp, Inc.

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CONNECTION
////////////////////////////////////////////////////////////////////////////////

// CreateBitbucketServerConnection - Creates new bitbucket server connection.
func (c *Client) CreateBitbucketServerConnection(connection BitbucketServerConnection) (*BitbucketServerConnection, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections", c.HostURL)
	return create(c, url, connection)
}

// ReadBitbucketServerConnection - Returns bitbucket server connection.
func (c *Client) ReadBitbucketServerConnection(id string) (*BitbucketServerConnection, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s", c.HostURL, id)
	return read[BitbucketServerConnection](c, url)
}

// UpdateBitbucketServerConnection - Updates bitbucket server connection.
func (c *Client) UpdateBitbucketServerConnection(id string, connection BitbucketServerConnection) (*BitbucketServerConnection, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s", c.HostURL, id)
	return update(c, url, connection)
}

// DeleteBitbucketServerConnection - Deletes a bitbucket server connection.
func (c *Client) DeleteBitbucketServerConnection(id string) error {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s", c.HostURL, id)
	return delete(c, url)
}

////////////////////////////////////////////////////////////////////////////////
// SCOPE CONFIG
////////////////////////////////////////////////////////////////////////////////

// CreateBitbucketServerConnectionScopeConfig - Creates a bitbucket server connection scope config.
func (c *Client) CreateBitbucketServerConnectionScopeConfig(connectionId string, scopeConfig BitbucketServerConnectionScopeConfig) (*BitbucketServerConnectionScopeConfig, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs", c.HostURL, connectionId)
	return create(c, url, scopeConfig)
}

// ReadBitbucketServerConnectionScopeConfig - Reads a bitbucket server connection scope config.
func (c *Client) ReadBitbucketServerConnectionScopeConfig(connectionId, scopeConfigId string) (*BitbucketServerConnectionScopeConfig, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId)
	return read[BitbucketServerConnectionScopeConfig](c, url)
}

// UpdateBitbucketServerConnectionScopeConfig - Updates a bitbucket server connection scope config.
func (c *Client) UpdateBitbucketServerConnectionScopeConfig(connectionId, scopeConfigId string, scopeConfig BitbucketServerConnectionScopeConfig) (*BitbucketServerConnectionScopeConfig, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId)
	return update(c, url, scopeConfig)
}

// DeleteBitbucketServerConnectionScopeConfig - Deletes a bitbucket server connection scope config.
func (c *Client) DeleteBitbucketServerConnectionScopeConfig(connectionId, scopeConfigId string) error {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId)
	return delete(c, url)
}

////////////////////////////////////////////////////////////////////////////////
// SCOPE
////////////////////////////////////////////////////////////////////////////////

// CreateBitbucketServerConnectionScope - Creates a bitbucket server connection scope.
func (c *Client) CreateBitbucketServerConnectionScope(connectionId string, scope BitbucketServerConnectionScope) (*BitbucketServerConnectionScope, error) {
	data := struct {
		Data []BitbucketServerConnectionScope `json:"data"`
	}{
		Data: []BitbucketServerConnectionScope{scope},
	}
	rb, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scopes", c.HostURL, connectionId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// endpoint accepts list but we only ever create one scope at a time, can
	// not use generic function because of this
	createdScopes := []BitbucketServerConnectionScope{}
	err = json.Unmarshal(body, &createdScopes)
	if err != nil {
		return nil, err
	}

	return &createdScopes[0], nil
}

// ReadBitbucketServerConnectionScope - Reads a bitbucket server connection scope.
func (c *Client) ReadBitbucketServerConnectionScope(connectionId, scopeId string) (*BitbucketServerConnectionScope, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// can not use generic read since response also contains the scope config
	res := struct {
		Scope       BitbucketServerConnectionScope       `json:"scope"`
		ScopeConfig BitbucketServerConnectionScopeConfig `json:"scopeConfig"`
	}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res.Scope, nil
}

// UpdateBitbucketServerConnectionScope - Updates a bitbucket server connection scope.
func (c *Client) UpdateBitbucketServerConnectionScope(connectionId, scopeId string, scopeConfig BitbucketServerConnectionScope) (*BitbucketServerConnectionScope, error) {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId)
	return update(c, url, scopeConfig)
}

// DeleteBitbucketServerConnectionScope - Deletes a bitbucket server connection scope.
func (c *Client) DeleteBitbucketServerConnectionScope(connectionId, scopeId string) error {
	url := fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId)
	return delete(c, url)
}
