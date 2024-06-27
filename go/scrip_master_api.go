// Package api provides the Go implementation for Scrip Master API operations.
package api

import (
	"encoding/json"
	"net/http"
)

// ScripMasterService handles communication with the scrip master related methods of the API.
type ScripMasterService struct {
	client *Client
}

// NewScripMasterService creates a new ScripMasterService.
func NewScripMasterService(client *Client) *ScripMasterService {
	return &ScripMasterService{client: client}
}

// ScripMasterRequest represents a request to the scrip master API.
type ScripMasterRequest struct {
	ExchangeSegment string `json:"exchangeSegment"`
}

// ScripMasterResponse represents a response from the scrip master API.
type ScripMasterResponse struct {
	Scrips []Scrip `json:"scrips"`
}

// Scrip represents a single scrip returned by the scrip master API.
type Scrip struct {
	Symbol          string `json:"symbol"`
	ScripName       string `json:"scripName"`
	ExchangeSegment string `json:"exchangeSegment"`
}

// GetScripMaster fetches the scrip master list from the API.
func (s *ScripMasterService) GetScripMaster(req *ScripMasterRequest) (*ScripMasterResponse, error) {
	url := s.client.BaseURL + "/scrip_master"
	resp, err := s.client.NewRequest(http.MethodGet, url, req)
	if err != nil {
		return nil, err
	}

	var scripMasterResp ScripMasterResponse
	err = json.NewDecoder(resp.Body).Decode(&scripMasterResp)
	if err != nil {
		return nil, err
	}

	return &scripMasterResp, nil
}
