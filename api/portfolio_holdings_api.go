// Package api provides the Go implementation for portfolio holdings API operations.
package go

import (
	"encoding/json"
	"net/http"
)

// PortfolioHoldingsService handles communication with the portfolio holdings related methods of the API.
type PortfolioHoldingsService service

// PortfolioHoldingsRequest represents a request for fetching portfolio holdings.
type PortfolioHoldingsRequest struct {
	ClientID string `json:"client_id"`
}

// PortfolioHoldingsResponse represents the response from the portfolio holdings API.
type PortfolioHoldingsResponse struct {
	Holdings []Holding `json:"holdings"`
}

// Holding represents a single holding in the portfolio.
type Holding struct {
	Symbol    string  `json:"symbol"`
	Quantity  int     `json:"quantity"`
	AveragePrice float64 `json:"average_price"`
	CurrentPrice float64 `json:"current_price"`
	ProfitLoss  float64 `json:"profit_loss"`
}

// Fetch retrieves the portfolio holdings for the specified client ID.
func (s *PortfolioHoldingsService) Fetch(req *PortfolioHoldingsRequest) (*PortfolioHoldingsResponse, *http.Response, error) {
	url := s.client.BaseURL + "/portfolio/holdings"
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	var holdingsResp PortfolioHoldingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&holdingsResp); err != nil {
		return nil, resp, err
	}

	return &holdingsResp, resp, nil
}
