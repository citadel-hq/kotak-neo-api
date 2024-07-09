// Package api provides the Go implementation for login API operations.
package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// LoginAPI is a client to perform requests to the login API.
type LoginAPI struct {
	APIClient   *APIClient
	Base64Token string
	RestClient  *RESTClientObject
}

// NewLoginAPI creates a new instance of LoginAPI.
func NewLoginAPI(apiClient *APIClient) *LoginAPI {
	return &LoginAPI{
		APIClient:   apiClient,
		Base64Token: apiClient.Config.Base64Token,
		RestClient:  apiClient.RestClient,
	}
}

// SessionInit initializes a session by sending a POST request to the specified URL with OAuth2 token.
func (api *LoginAPI) SessionInit() (map[string]interface{}, error) {
	headerParams := map[string]string{"Authorization": "Basic " + api.Base64Token}
	bodyParams := map[string]string{"grant_type": "client_credentials"}
	URL, err := api.APIClient.Config.getDomain(true)
	URL = URL + "oauth2/token"

	response, err := api.RestClient.Request(http.MethodPost, URL, nil, headerParams, bodyParams)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		var jsonResp map[string]interface{}
		json.NewDecoder(response.Body).Decode(&jsonResp)
		api.APIClient.Config.BearerToken = jsonResp["access_token"].(string)
		return jsonResp, nil
	}
	return map[string]interface{}{"data": map[string]interface{}{"Code": response.StatusCode, "Message": "Error occurred to initialise the session"}}, nil
}

// GenerateViewToken generates a view token for a given mobile number and password.
func (api *LoginAPI) GenerateViewToken(password, mobilenumber, userid, pan, mpin string) (map[string]interface{}, error) {
	headerParams := map[string]string{"Authorization": "Bearer " + api.APIClient.Config.BearerToken}

	bodyParams, err := LoginParamsValidation(mobilenumber, userid, pan, password, mpin)
	if err != nil {
		return nil, err
	}

	api.APIClient.Config.LoginParams = bodyParams
	URL, err := api.APIClient.Config.getUrlDetails("view_token")
	if err != nil {
		return nil, err
	}

	response, err := api.RestClient.Request(http.MethodPost, URL, nil, headerParams, bodyParams)
	if err != nil {
		return nil, err
	}

	var jsonResp map[string]interface{}
	json.NewDecoder(response.Body).Decode(&jsonResp)

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		if mobilenumber != "" && !strings.HasPrefix(mobilenumber, "+") {
			jsonResp["message"] = "since no country code found we have appended +91 as the default country code. Please change it to the correct code if your mobile number is not of Indian number."
		}
		data := jsonResp["data"].(map[string]interface{})
		api.APIClient.Config.ViewToken = data["token"].(string)
		api.APIClient.Config.Sid = data["sid"].(string)
		return jsonResp, nil
	}

	if mobilenumber != "" && !strings.HasPrefix(mobilenumber, "+") {
		jsonResp["Note"] = "since no country code found we have appended +91 as the default country code. Please change it to the correct code if your mobile number is not of Indian number."
	}
	return jsonResp, nil
}

// GenerateOTP generates an OTP for the user.
func (api *LoginAPI) GenerateOTP() (string, error) {
	headerParams := map[string]string{"Authorization": "Bearer " + api.APIClient.Config.BearerToken}
	userId, err := api.APIClient.Config.extractUserId(api.APIClient.Config.ViewToken)
	if err != nil {
		return "", err
	}
	bodyParams := map[string]interface{}{
		"userId":        userId,
		"sendEmail":     true,
		"isWhitelisted": true,
	}
	URL, err := api.APIClient.Config.getUrlDetails("generate_otp")
	if err != nil {
		return "", err
	}

	response, err := api.RestClient.Request(http.MethodPost, URL, nil, headerParams, bodyParams)
	if err != nil {
		return "", err
	}
	var jsonResp map[string]interface{}
	json.NewDecoder(response.Body).Decode(&jsonResp)

	// TODO: Recheck this?
	return jsonResp["text"].(string), nil
}

// Login2FA performs two-factor authentication for login.
func (api *LoginAPI) Login2FA(OTP string) (map[string]interface{}, error) {
	params := api.APIClient.Config.LoginParams
	bodyParams := map[string]string{}

	if mobileNumber, ok := params["mobileNumber"]; ok && len(OTP) == 6 {
		bodyParams["mobileNumber"] = mobileNumber
		bodyParams["mpin"] = OTP
	} else if pan, ok := params["pan"]; ok && len(OTP) == 6 {
		bodyParams["pan"] = pan
		bodyParams["mpin"] = OTP
	} else {
		bodyParams["userId"] = api.APIClient.Config.UserId
		bodyParams["otp"] = OTP
	}

	headerParams := map[string]string{
		"Authorization": "Bearer " + api.APIClient.Config.BearerToken,
		"sid":           api.APIClient.Config.Sid,
		"Auth":          api.APIClient.Config.ViewToken,
	}

	URL, err := api.APIClient.Config.getUrlDetails("edit_token")
	if err != nil {
		return nil, err
	}

	response, err := api.RestClient.Request(http.MethodPost, URL, nil, headerParams, bodyParams)
	if err != nil {
		return nil, err
	}

	var jsonResp map[string]interface{}
	json.NewDecoder(response.Body).Decode(&jsonResp)

	if _, ok := jsonResp["error"]; !ok {
		data := jsonResp["data"].(map[string]interface{})
		api.APIClient.Config.EditToken = data["token"].(string)
		api.APIClient.Config.EditSid = data["sid"].(string)
		api.APIClient.Config.EditRid = data["rid"].(string)
		api.APIClient.Config.ServerId = data["hsServerId"].(string)
	}
	return jsonResp, nil
}
