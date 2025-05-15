// Copyright (c) HashiCorp, Inc.

package client

type ApiKey struct {
	ID           int    `json:"id"`
	AllowedPath  string `json:"allowedPath"`
	ApiKey       string `json:"apiKey"`
	CreatedAt    string `json:"createdAt"`
	Creator      string `json:"creator"`
	CreatorEmail string `json:"creatorEmail"`
	ExpiredAt    string `json:"expiredAt"`
	Extra        string `json:"extra"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	UpdatedAt    string `json:"updatedAt"`
	Updater      string `json:"updater"`
	UpdaterEmail string `json:"updaterEmail"`
}

type ApiKeyCreate struct {
	AllowedPath string `json:"allowedPath"`
	ExpiredAt   string `json:"expiredAt"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type JsonBody struct {
	Causes  []string `json:"causes"`
	Code    int      `json:"code"`
	Data    string   `json:"data"`
	Message string   `json:"message"`
	Success bool     `json:"success"`
}

type BitbucketServerConnection struct {
	ID               int    `json:"id"`
	CreatedAt        string `json:"createdAt"`
	Endpoint         string `json:"endpoint"`
	Name             string `json:"name"`
	Password         string `json:"password"`
	Proxy            string `json:"proxy"`
	RateLimitPerHour int    `json:"rateLimitPerHour"`
	UpdatedAt        string `json:"updatedAt"`
	Username         string `json:"username"`
}

type BitbucketServerConnectionScopeConfig struct {
	ConnectionId int      `json:"connectionId"`
	CreatedAt    string   `json:"createdAt"`
	ID           int      `json:"id"`
	Entities     []string `json:"entities"`
	Name         string   `json:"name"`
	PrComponent  string   `json:"prComponent"`
	PrType       string   `json:"prType"`
	RefDiff      *RefDiff `json:"refdiff"`
	UpdatedAt    string   `json:"updatedAt"`
}

type RefDiff struct {
	TagsLimit   int    `json:"tagsLimit"`
	TagsPattern string `json:"tagsPattern"`
}

type BitbucketServerConnectionScope struct {
	BitbucketId   string `json:"bitbucketId"`
	CloneUrl      string `json:"cloneUrl"`
	ConnectionId  int    `json:"connectionId"`
	CreatedAt     string `json:"createdAt"`
	Description   string `json:"description"`
	HTMLUrl       string `json:"HTMLUrl"`
	Name          string `json:"name"`
	ScopeConfigId int    `json:"scopeConfigId"`
	UpdatedAt     string `json:"updatedAt"`
}

type GithubConnection struct {
	ID               int    `json:"id"`
	AppId            string `json:"appId"`
	AuthMethod       string `json:"authMethod"`
	CreatedAt        string `json:"createdAt"`
	EnableGraphql    bool   `json:"enableGraphql"`
	Endpoint         string `json:"endpoint"`
	InstallationId   int    `json:"installationId"`
	Name             string `json:"name"`
	Proxy            string `json:"proxy"`
	RateLimitPerHour int    `json:"rateLimitPerHour"`
	SecretKey        string `json:"secretKey"`
	Token            string `json:"token"`
	UpdatedAt        string `json:"updatedAt"`
}

type GithubConnectionScopeConfig struct {
	ConnectionId int      `json:"connectionId"`
	CreatedAt    string   `json:"createdAt"`
	ID           int      `json:"id"`
	Entities     []string `json:"entities"`
	Name         string   `json:"name"`
	PrComponent  string   `json:"prComponent"`
	PrType       string   `json:"prType"`
	RefDiff      *RefDiff `json:"refdiff"`
	UpdatedAt    string   `json:"updatedAt"`
}

type GithubConnectionScope struct {
	BitbucketId   string `json:"bitbucketId"`
	CloneUrl      string `json:"cloneUrl"`
	ConnectionId  int    `json:"connectionId"`
	CreatedAt     string `json:"createdAt"`
	Description   string `json:"description"`
	HTMLUrl       string `json:"HTMLUrl"`
	Name          string `json:"name"`
	ScopeConfigId int    `json:"scopeConfigId"`
	UpdatedAt     string `json:"updatedAt"`
}
