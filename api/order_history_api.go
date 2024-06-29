// Package api implements the order history API in Go.
package api

import (
	"fmt"
	"encoding/json"
	"net/http"
)

// OrderHistoryService provides operations to fetch order history.
type OrderHistoryService struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewOrderHistoryService creates a new service for order history.
func NewOrderHistoryService(apiClient *APIClient) *OrderHistoryService {
	return &OrderHistoryService{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// FetchOrderHistory fetches the history of an order.
func (api *OrderHistoryService) FetchOrderHistory(orderID string) (map[string]interface{}, error) {
	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           api.ApiClient.Config.EditSid,
		"Auth":          api.ApiClient.Config.EditToken,
		"neo-fin-key":   api.ApiClient.Config.getNeoFinKey(),
		"accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	queryParams := map[string]string{
		"sId": api.ApiClient.Config.ServerId,
	}
	bodyParams := map[string]string{
		"nOrdNo": orderID,
	}
	url, err := api.ApiClient.Config.getUrlDetails("order_history")
	if err != nil {
		return nil, err
	}
	resp, err := api.RestClient.Request(http.MethodPost, url, queryParams, headerParams, bodyParams)
	if err != nil {
		return nil, err
	}
	var jsonResp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)
	return jsonResp, nil
}
