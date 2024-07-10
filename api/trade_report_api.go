// Package api provides the Go implementation for trade report API operations.
package api

import (
	"encoding/json"
	"net/http"
	"fmt"
	"errors"
)

// TradeReportAPI struct defines the client for the trade report API operations.
type TradeReportAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewTradeReportAPI creates a new trade report API client with the specified base URL.
func NewTradeReportAPI(apiClient *APIClient) *TradeReportAPI {
	return &TradeReportAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// GetTradeReport fetches the trade report for all trades or a specific trade by order ID.
func (api *TradeReportAPI) GetTradeReport(orderID string) (map[string]interface{}, error) {
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
	url, err := api.ApiClient.Config.getUrlDetails("trade_report")
	if err != nil {
		return nil, err
	}
	resp, err := api.RestClient.Request(http.MethodGet, url, queryParams, headerParams, nil)
	if err != nil {
		return nil, err
	}
	var tradeReport map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&tradeReport)
	if orderID != "" {
		outputJson := map[string]interface{}{}
		if data, ok := tradeReport["data"].([]map[string]interface{}); ok {
			outputJson["tid"] = tradeReport["tid"].(interface{})
			outputJson["stat"] = tradeReport["stat"].(interface{})
			outputJson["stCode"] = tradeReport["stCode"].(interface{})
			for _, item := range data {
				if orderNumber, ok := item["nOrdNo"].(string); ok {
					if orderNumber == orderID {
						outputJson["data"] = item
					}
				}
			}
			return outputJson, nil
		}
		return nil, errors.New("there is no trades available with the given order id")
	}
	return nil, errors.New("no orderId sent in trade report")
}
