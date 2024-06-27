// Package api_client provides the Go implementation for API client operations.
package api

const HeaderUserAgent string = "User-Agent"

// APIClient struct defines the client for the API operations.
type APIClient struct {
	RestClient     *RESTClientObject
	UserAgent      string
	DefaultHeaders map[string]string
	Config         Configuration
}

// NewAPIClient creates a new API client with the specified base URL.
func NewAPIClient(configuration Configuration) *APIClient {
	return &APIClient{
		RestClient:     NewRESTClientObject(configuration),
		UserAgent:      "NeoTradeApi-python/1.0.0/python",
		DefaultHeaders: map[string]string{},
		Config:         configuration,
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
