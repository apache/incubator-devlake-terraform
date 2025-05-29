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

// CreateGithubConnection - Creates new github connection.
func (c *Client) CreateGithubConnection(connection GithubConnection) (*GithubConnection, error) {
	url := fmt.Sprintf("%s/plugins/github/connections", c.HostURL)
	return create(c, url, connection)
}

// ReadGithubConnection - Returns github connection.
func (c *Client) ReadGithubConnection(id string) (*GithubConnection, error) {
	url := fmt.Sprintf("%s/plugins/github/connections/%s", c.HostURL, id)
	return read[GithubConnection](c, url)
}

// UpdateGithubConnection - Updates github connection.
func (c *Client) UpdateGithubConnection(id string, connection GithubConnection) (*GithubConnection, error) {
	url := fmt.Sprintf("%s/plugins/github/connections/%s", c.HostURL, id)
	return update(c, url, connection)
}

// DeleteGithubConnection - Deletes a github connection.
func (c *Client) DeleteGithubConnection(id string) error {
	url := fmt.Sprintf("%s/plugins/github/connections/%s", c.HostURL, id)
	return del(c, url)
}

////////////////////////////////////////////////////////////////////////////////
// SCOPE CONFIG
////////////////////////////////////////////////////////////////////////////////

// CreateGithubConnectionScopeConfig - Creates a github connection scope config.
func (c *Client) CreateGithubConnectionScopeConfig(connectionId string, scopeConfig GithubConnectionScopeConfig) (*GithubConnectionScopeConfig, error) {
	url := fmt.Sprintf("%s/plugins/github/connections/%s/scope-configs", c.HostURL, connectionId)
	return create(c, url, scopeConfig)
}

// ReadGithubConnectionScopeConfig - Reads a github connection scope config.
func (c *Client) ReadGithubConnectionScopeConfig(connectionId, scopeConfigId string) (*GithubConnectionScopeConfig, error) {
	url := fmt.Sprintf("%s/plugins/github/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId)
	return read[GithubConnectionScopeConfig](c, url)
}

// UpdateGithubConnectionScopeConfig - Updates a github connection scope config.
func (c *Client) UpdateGithubConnectionScopeConfig(connectionId, scopeConfigId string, scopeConfig GithubConnectionScopeConfig) (*GithubConnectionScopeConfig, error) {
	url := fmt.Sprintf("%s/plugins/github/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId)
	return update(c, url, scopeConfig)
}

// DeleteGithubConnectionScopeConfig - Deletes a github connection scope config.
func (c *Client) DeleteGithubConnectionScopeConfig(connectionId, scopeConfigId string) error {
	url := fmt.Sprintf("%s/plugins/github/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId)
	return del(c, url)
}

////////////////////////////////////////////////////////////////////////////////
// SCOPE
////////////////////////////////////////////////////////////////////////////////

// CreateGithubConnectionScope - Creates a github connection scope.
func (c *Client) CreateGithubConnectionScope(connectionId string, scope GithubConnectionScope) (*GithubConnectionScope, error) {
	data := struct {
		Data []GithubConnectionScope `json:"data"`
	}{
		Data: []GithubConnectionScope{scope},
	}
	rb, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/plugins/github/connections/%s/scopes", c.HostURL, connectionId), strings.NewReader(string(rb)))
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
	createdScopes := []GithubConnectionScope{}
	err = json.Unmarshal(body, &createdScopes)
	if err != nil {
		return nil, err
	}

	return &createdScopes[0], nil
}

// ReadGithubConnectionScope - Reads a github connection scope.
func (c *Client) ReadGithubConnectionScope(connectionId, scopeId string) (*GithubConnectionScope, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/plugins/github/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId), nil)
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
		Scope       GithubConnectionScope       `json:"scope"`
		ScopeConfig GithubConnectionScopeConfig `json:"scopeConfig"`
	}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res.Scope, nil
}

// UpdateGithubConnectionScope - Updates a github connection scope.
func (c *Client) UpdateGithubConnectionScope(connectionId, scopeId string, scopeConfig GithubConnectionScope) (*GithubConnectionScope, error) {
	url := fmt.Sprintf("%s/plugins/github/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId)
	return update(c, url, scopeConfig)
}

// DeleteGithubConnectionScope - Deletes a github connection scope.
func (c *Client) DeleteGithubConnectionScope(connectionId, scopeId string) error {
	url := fmt.Sprintf("%s/plugins/github/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId)
	return del(c, url)
}
