package api

import (
	"fmt"
	"encoding/json"
	"net/http"
)

// OrderReportAPI handles the order report functionality
type OrderReportAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewOrderReportAPI creates a new OrderReportAPI instance
func NewOrderReportAPI(apiClient *APIClient) *OrderReportAPI {
	return &OrderReportAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// GetOrderReport fetches the order report for the given order ID
func (api *OrderReportAPI) GetOrderReport() (map[string]interface{}, error) {
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

	url, err := api.ApiClient.Config.getUrlDetails("order_book")
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
