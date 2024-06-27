package go

import (
	"encoding/json"
	"net/http"
)

// PositionsAPI represents the structure for positions API operations
type PositionsAPI struct {
	Client *http.Client
}

// NewPositionsAPI creates a new instance of PositionsAPI with the provided HTTP client
func NewPositionsAPI(client *http.Client) *PositionsAPI {
	return &PositionsAPI{
		Client: client,
	}
}

// GetPositions retrieves the positions for the user
func (api *PositionsAPI) GetPositions() ([]Position, error) {
	// Define the URL for the positions API endpoint
	url := "http://example.com/api/positions"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response body into the positions slice
	var positions []Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

// Position represents the structure of a trading position
type Position struct {
	ID          string  `json:"id"`
	Instrument  string  `json:"instrument"`
	Quantity    int     `json:"quantity"`
	AverageCost float64 `json:"averageCost"`
	CurrentPrice float64 `json:"currentPrice"`
	ProfitLoss  float64 `json:"profitLoss"`
}
