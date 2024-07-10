package api

import (
	"encoding/json"
	"net/http"
	"fmt"
)

// PositionsAPI represents the structure for positions API operations
type PositionsAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewPositionsAPI creates a new instance of PositionsAPI with the provided HTTP client
func NewPositionsAPI(apiClient *APIClient) *PositionsAPI {
	return &PositionsAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// GetPositions retrieves the positions for the user
func (api *PositionsAPI) GetPositions() (map[string]interface{}, error) {
	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           api.ApiClient.Config.EditSid,
		"Auth":          api.ApiClient.Config.EditToken,
		"neo-fin-key":   api.ApiClient.Config.getNeoFinKey(),
		"accept":        "application/json",
	}
	queryParams := map[string]string{
		"sId": api.ApiClient.Config.ServerId,
	}
	url, err := api.ApiClient.Config.getUrlDetails("positions")
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
