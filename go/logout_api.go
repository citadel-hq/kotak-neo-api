// Package api provides the Go implementation for logout API operations.
package api

import (
	"net/http"
)

// LogoutAPI struct defines the client for the logout API operations.
type LogoutAPI struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewLogoutAPI creates a new logout API client with the specified base URL.
func NewLogoutAPI(baseURL string) *LogoutAPI {
	return &LogoutAPI{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// LogoutUser performs a logout operation for a user session.
func (api *LogoutAPI) LogoutUser() error {
	url := api.BaseURL + "/logout"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to logout user, status code: %d", resp.StatusCode)
	}

	return nil
}
