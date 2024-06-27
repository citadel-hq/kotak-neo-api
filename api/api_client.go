// Package api_client provides the Go implementation for API client operations.
package api

import (
	"net/http"
)

const HeaderUserAgent string = "User-Agent"

// APIClient struct defines the client for the API operations.
type APIClient struct {
	BaseURL        string
	HTTPClient     *http.Client
	UserAgent      string
	DefaultHeaders map[string]string
}

// NewAPIClient creates a new API client with the specified base URL.
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL:        baseURL,
		HTTPClient:     &http.Client{},
		UserAgent:      "NeoTradeApi-python/1.0.0/python",
		DefaultHeaders: map[string]string{},
	}
}

// Implement API client methods here.

func (c *APIClient) GetUserAgent() string {
	return c.DefaultHeaders[HeaderUserAgent]
}

func (c *APIClient) SetUserAgent(agent string) {
	c.DefaultHeaders[HeaderUserAgent] = agent
}

func (c *APIClient) SetHeader(key string, val string) {
	c.DefaultHeaders[key] = val
}
