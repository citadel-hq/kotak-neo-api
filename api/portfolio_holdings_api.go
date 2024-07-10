// Package api provides the Go implementation for portfolio holdings API operations.
package api

import (
	"encoding/json"
	"net/http"
	"fmt"
)

type PortfolioHoldingsAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

func NewPortfolioHoldingsAPI(apiClient *APIClient) *PortfolioHoldingsAPI {
	return &PortfolioHoldingsAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// Holding represents a single holding in the portfolio.
type Holding struct {
	Symbol       string  `json:"symbol"`
	Quantity     int     `json:"quantity"`
	AveragePrice float64 `json:"average_price"`
	CurrentPrice float64 `json:"current_price"`
	ProfitLoss   float64 `json:"profit_loss"`
}

// Fetch retrieves the portfolio holdings for the specified client ID.
func (api *PortfolioHoldingsAPI) Fetch() (map[string]interface{}, error) {
	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           api.ApiClient.Config.EditSid,
		"Auth":          api.ApiClient.Config.EditToken,
		"accept":        "*/*",
	}
	queryParams := map[string]string{
		"sId": api.ApiClient.Config.ServerId,
	}
	url, err := api.ApiClient.Config.getUrlDetails("holdings")
	if err != nil {
		return nil, err
	}
	resp, err := api.RestClient.Request(http.MethodGet, url, queryParams, headerParams, nil)
	if err != nil {
		return nil, err
	}
	var jsonResp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)
	return jsonResp, nil
}
