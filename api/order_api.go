// Package api implements the order API in Go.
package api

import (
	"fmt"
	"encoding/json"
	"net/http"
)

// OrderService provides operations to manage orders.
type OrderService struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewOrderService creates a new service for orders.
func NewOrderService(apiClient *APIClient) *OrderService {
	return &OrderService{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// PlaceOrder places a new order.
func (api *OrderService) PlaceOrder(request OrderRequest) (map[string]interface{}, error) {
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
	var bodyParams map[string]interface{}
	jsonString, _ := json.Marshal(request)
	_ = json.Unmarshal(jsonString, &bodyParams)
	url, err := api.ApiClient.Config.getUrlDetails("place_order")
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

// CancelOrder cancels an existing order.
func (api *OrderService) CancelOrder(orderID string, isVerify bool, amo float64) (map[string]interface{}, error) {
	if isVerify {
		orderBookResp, err := (&OrderReportAPI{ApiClient: api.ApiClient}).GetOrderReport()
		if err != nil {
			return nil, err
		}

		if data, ok := orderBookResp["data"].([]interface{}); ok {
			for _, item := range data {
				if order, ok := item.(map[string]interface{}); ok && order["nOrdNo"] == orderID {
					status := order["ordSt"].(string)
					if status == "rejected" || status == "cancelled" || status == "complete" || status == "traded" {
						if status == "complete" {
							status = "Traded"
						}
						return map[string]interface{}{
							"Error":  fmt.Sprintf("The Given Order Status is %s", status),
							"Reason": order["rejRsn"],
						}, nil
					}
				}
			}
		}
	}

	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           api.ApiClient.Config.EditSid,
		"Auth":          api.ApiClient.Config.EditToken,
		"neo-fin-key":   api.ApiClient.Config.getNeoFinKey(),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	bodyParams := map[string]interface{}{
		"on": orderID,
		"am": amo,
	}

	queryParams := map[string]string{
		"sId": api.ApiClient.Config.ServerId,
	}

	url, err := api.ApiClient.Config.getUrlDetails("cancel_order")
	if err != nil {
		return nil, err
	}

	resp, err := api.RestClient.Request(http.MethodPost, url, queryParams, headerParams, bodyParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jsonResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, err
	}

	return jsonResp, nil
}

type OrderRequest struct {
	Amo               float64 `json:"am,omitempty"`
	DisclosedQuantity int     `json:"dq,omitempty"`
	ExchangeSegment   string  `json:"es,omitempty"`
	MarketProtection  string  `json:"mp,omitempty"`
	Product           string  `json:"pc,omitempty"`
	Pf                string  `json:"pf,omitempty"`
	Price             float64 `json:"pr,omitempty"`
	OrderType         string  `json:"pt,omitempty"`
	Quantity          int     `json:"qt,omitempty"`
	Validity          string  `json:"rt,omitempty"`
	TriggerPrice      float64 `json:"tp,omitempty"`
	TradingSymbol     string  `json:"ts,omitempty"`
	TransactionType   string  `json:"tt,omitempty"`
	Tag               string  `json:"ig,omitempty"`
}

// OrderResponse represents a response from an order operation.
type OrderResponse struct {
	// Add order response fields here
}
