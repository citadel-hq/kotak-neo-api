// Package api provides the Go implementation for logout API operations.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LogoutAPI struct defines the client for the logout API operations.
type LogoutAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewLogoutAPI creates a new logout API client with the specified base URL.
func NewLogoutAPI(apiClient *APIClient) *LogoutAPI {
	return &LogoutAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// LogoutUser performs a logout operation for a user session.
func (api *LogoutAPI) LogoutUser() (map[string]interface{}, error) {
	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           api.ApiClient.Config.EditSid,
		"Auth":          api.ApiClient.Config.EditToken,
		"accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	url, err := api.ApiClient.Config.getUrlDetails("logout")
	if err != nil {
		return nil, err
	}
	resp, err := api.RestClient.Request(http.MethodPost, url, nil, headerParams, nil)
	if err != nil {
		return nil, err
	}
	var jsonResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&jsonResp)

	return jsonResp, nil
}
