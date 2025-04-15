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

// ReadBitbucketServerConnection - Returns bitbucket server connection.
func (c *Client) ReadBitbucketServerConnection(id string) (*BitbucketServerConnection, error) {
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

// DeleteBitbucketServerConnection - Deletes a bitbucket server connection.
func (c *Client) DeleteBitbucketServerConnection(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SCOPE CONFIG
////////////////////////////////////////////////////////////////////////////////

// CreateBitbucketServerConnectionScopeConfig - Creates a bitbucket server connection scope config.
func (c *Client) CreateBitbucketServerConnectionScopeConfig(connectionId string, scopeConfig BitbucketServerConnectionScopeConfig) (*BitbucketServerConnectionScopeConfig, error) {
	rb, err := json.Marshal(scopeConfig)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs", c.HostURL, connectionId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	createdScopeConfig := BitbucketServerConnectionScopeConfig{}
	err = json.Unmarshal(body, &createdScopeConfig)
	if err != nil {
		return nil, err
	}

	return &createdScopeConfig, nil
}

// ReadBitbucketServerConnectionScopeConfig - Reads a bitbucket server connection scope config.
func (c *Client) ReadBitbucketServerConnectionScopeConfig(connectionId, scopeConfigId string) (*BitbucketServerConnectionScopeConfig, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	scopeConfig := BitbucketServerConnectionScopeConfig{}
	err = json.Unmarshal(body, &scopeConfig)
	if err != nil {
		return nil, err
	}

	return &scopeConfig, nil
}

// UpdateBitbucketServerConnectionScopeConfig - Updates a bitbucket server connection scope config.
func (c *Client) UpdateBitbucketServerConnectionScopeConfig(connectionId, scopeConfigId string, scopeConfig BitbucketServerConnectionScopeConfig) (*BitbucketServerConnectionScopeConfig, error) {
	rb, err := json.Marshal(scopeConfig)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedScopeConfig := BitbucketServerConnectionScopeConfig{}
	err = json.Unmarshal(body, &updatedScopeConfig)
	if err != nil {
		return nil, err
	}

	return &updatedScopeConfig, nil
}

// DeleteBitbucketServerConnectionScopeConfig - Deletes a bitbucket server connection scope config.
func (c *Client) DeleteBitbucketServerConnectionScopeConfig(connectionId, scopeConfigId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scope-configs/%s", c.HostURL, connectionId, scopeConfigId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
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

	// endpoint accepts list but we only ever create one scope at a time
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
	rb, err := json.Marshal(scopeConfig)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedScope := BitbucketServerConnectionScope{}
	err = json.Unmarshal(body, &updatedScope)
	if err != nil {
		return nil, err
	}

	return &updatedScope, nil
}

// DeleteBitbucketServerConnectionScope - Deletes a bitbucket server connection scope.
func (c *Client) DeleteBitbucketServerConnectionScope(connectionId, scopeId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/plugins/bitbucket_server/connections/%s/scopes/%s", c.HostURL, connectionId, scopeId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
