// Package api provides the Go implementation for Scrip Search API operations.
package api

import (
	"encoding/json"
	"net/http"
)

// ScripSearchService handles communication with the scrip search related methods of the API.
type ScripSearchService struct {
	client *Client
}

// NewScripSearchService creates a new ScripSearchService.
func NewScripSearchService(client *Client) *ScripSearchService {
	return &ScripSearchService{client: client}
}

// ScripSearchRequest represents a request to the scrip search API.
type ScripSearchRequest struct {
	Keyword string `json:"keyword"`
}

// ScripSearchResponse represents a response from the scrip search API.
type ScripSearchResponse struct {
	Scrips []Scrip `json:"scrips"`
}

// Scrip represents a single scrip returned by the scrip search API.
type Scrip struct {
	Symbol          string `json:"symbol"`
	ScripName       string `json:"scripName"`
	ExchangeSegment string `json:"exchangeSegment"`
}

// SearchScrips fetches the scrips based on the search keyword from the API.
func (s *ScripSearchService) SearchScrips(req *ScripSearchRequest) (*ScripSearchResponse, error) {
	url := s.client.BaseURL + "/scrip_search"
	resp, err := s.client.NewRequest(http.MethodGet, url, req)
	if err != nil {
		return nil, err
	}

	var scripSearchResp ScripSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&scripSearchResp)
	if err != nil {
		return nil, err
	}

	return &scripSearchResp, nil
}
