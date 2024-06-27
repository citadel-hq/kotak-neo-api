// Package api provides the Go implementation for login API operations.
package api

import (
	"encoding/json"
	"net/http"
	"time"
	"bytes"
)

// LoginAPI struct defines the client for the login API operations.
type LoginAPI struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewLoginAPI creates a new login API client with the specified base URL.
func NewLoginAPI(baseURL string) *LoginAPI {
	return &LoginAPI{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SessionInit initializes a session by sending a POST request with OAuth2 token.
func (api *LoginAPI) SessionInit() (*json.RawMessage, error) {
	url := api.BaseURL + "/oauth2/token"
	reqBody := map[string]string{
		"grant_type": "client_credentials",
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := api.HTTPClient.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateViewToken generates a view token for a given mobile number and password.
func (api *LoginAPI) GenerateViewToken(mobileNumber, password string) (*json.RawMessage, error) {
	url := api.BaseURL + "/view_token"
	reqBody := map[string]string{
		"mobileNumber": mobileNumber,
		"password":     password,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := api.HTTPClient.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Login2FA handles the second factor authentication for login.
func (api *LoginAPI) Login2FA(otp string) (*json.RawMessage, error) {
	url := api.BaseURL + "/login_2fa"
	reqBody := map[string]string{
		"otp": otp,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := api.HTTPClient.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
