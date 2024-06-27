// Package api provides the Go implementation for limits API operations.
package api

import (
	"encoding/json"
	"net/http"
)

// LimitsAPI struct defines the client for the limits API operations.
type LimitsAPI struct {
	Client *http.Client
	BaseURL string
}

// NewLimitsAPI creates a new limits API client with the specified base URL.
func NewLimitsAPI(baseURL string) *LimitsAPI {
	return &LimitsAPI{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

// GetLimits retrieves the limits for a given segment, exchange, and product.
func (api *LimitsAPI) GetLimits(segment, exchange, product string) (map[string]interface{}, error) {
	url := api.BaseURL + "/limits"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("segment", segment)
	q.Add("exchange", exchange)
	q.Add("product", product)
	req.URL.RawQuery = q.Encode()

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
