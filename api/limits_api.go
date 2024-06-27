// Package api provides the Go implementation for limits API operations.
package api

import (
	"encoding/json"
	"fmt"
)

// LimitsAPI struct defines the client for the limits API operations.
type LimitsAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewLimitsAPI creates a new limits API client with the specified base URL.
func NewLimitsAPI(apiClient *APIClient) *LimitsAPI {
	return &LimitsAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// GetLimits retrieves the limits for a given segment, exchange, and product.
func (api *LimitsAPI) GetLimits(segment, exchange, product string) (map[string]interface{}, error) {
	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           fmt.Sprintf("Bearer %s", api.ApiClient.Config.EditSid),
		"Auth":          fmt.Sprintf("Bearer %s", api.ApiClient.Config.EditToken),
		"neo-fin-key":   fmt.Sprintf("Bearer %s", api.ApiClient.Config.getNeoFinKey()),
		"accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	queryParams := map[string]string{
		"sId": fmt.Sprintf("Bearer %s", api.ApiClient.Config.ServerId),
	}

	bodyParams := map[string]string{
		"seg":  segment,
		"exch": exchange,
		"prod": product,
	}

	url, err := api.ApiClient.Config.getUrlDetails("limits")
	if err != nil {
		return nil, err
	}

	limits_report, err := api.RestClient.Request(url, "POST", queryParams, headerParams, bodyParams)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(limits_report.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
