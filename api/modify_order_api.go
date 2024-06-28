// Package api provides the Go implementation for modifying orders.
package api

import (
	"encoding/json"
	"net/http"
	"fmt"
)

type ModifyOrderAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

func NewModifyOrderAPI(apiClient *APIClient) *ModifyOrderAPI {
	return &ModifyOrderAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

func (api *ModifyOrderAPI) ModifyQuick(modificationRequest ModifyOrderRequest) (map[string]interface{}, error) {
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
	jsonString, _ := json.Marshal(modificationRequest)
	_ = json.Unmarshal(jsonString, &bodyParams)

	url, err := api.ApiClient.Config.getUrlDetails("modify_order")
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

func (api *ModifyOrderAPI) ModifyWithOrderID(modificationRequest ModifyOrderRequest) (map[string]interface{}, error) {
	headerParams := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
		"Sid":           api.ApiClient.Config.EditSid,
		"Auth":          api.ApiClient.Config.EditToken,
		"neo-fin-key":   api.ApiClient.Config.getNeoFinKey(),
		"accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	reportAPI := OrderReportAPI{
		ApiClient:  api.ApiClient,
		RestClient: api.RestClient,
	}
	orderBookResp, err := reportAPI.GetOrderReport()
	if err != nil {
		return nil, err
	}
	if data, ok := orderBookResp["data"].([]interface{}); ok {
		for _, item := range data {
			order := item.(map[string]interface{})
			if order["nOrdNo"] == modificationRequest.OrderID {
				if orderStatus, ok := order["ordSt"].(string); ok && (orderStatus == "rejected" || orderStatus == "cancelled" || orderStatus == "complete" || orderStatus == "traded") {
					if orderStatus == "complete" {
						orderStatus = "Traded"
					}
					return map[string]interface{}{
						"Error":  fmt.Sprintf("The Given Order Status is %s, So we can't proceed further", orderStatus),
						"Reason": order["rejRsn"],
					}, nil
				}

				modificationRequest.TradingSymbol = defaultString(modificationRequest.TradingSymbol, order["trdSym"].(string))
				modificationRequest.InstrumentToken = defaultString(modificationRequest.InstrumentToken, order["tok"].(string))
				modificationRequest.Product = defaultString(modificationRequest.Product, order["prod"].(string))
				modificationRequest.TransactionType = defaultString(modificationRequest.TransactionType, order["trnsTp"].(string))
				modificationRequest.ExchangeSegment = defaultString(modificationRequest.ExchangeSegment, order["exSeg"].(string))
				var triggerPrice float64
				if modificationRequest.TriggerPrice != 0 {
					triggerPrice = modificationRequest.TriggerPrice
				} else {
					triggerPrice = order["trgPrc"].(float64)
				}
				modificationRequest.TriggerPrice = triggerPrice

				var bodyParams map[string]interface{}
				jsonString, _ := json.Marshal(modificationRequest)
				_ = json.Unmarshal(jsonString, &bodyParams)

				queryParams := map[string]string{
					"sId": api.ApiClient.Config.ServerId,
				}

				url, err := api.ApiClient.Config.getUrlDetails("modify_order")
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
		}
	}

	return map[string]interface{}{
		"Message": fmt.Sprintf("The Given Order Number is %s and it is not matching with any order of the orders", modificationRequest.OrderID),
	}, nil
}

func defaultString(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

// ModifyOrderRequest represents the request structure for modifying an order.
type ModifyOrderRequest struct {
	OrderID           string  `json:"no,omitempty"` // order_id
	Price             float64 `json:"pr,omitempty"` // price
	OrderType         string  `json:"pt,omitempty"` // order_type
	Quantity          int     `json:"qt,omitempty"` // quantity
	Validity          string  `json:"vd,omitempty"` // validity
	InstrumentToken   string  `json:"tk,omitempty"` // instrument_token
	ExchangeSegment   string  `json:"es,omitempty"` // exchange_segment
	Product           string  `json:"pc,omitempty"` // product
	TradingSymbol     string  `json:"ts,omitempty"` // trading_symbol
	TransactionType   string  `json:"tt,omitempty"` // transaction_type
	TriggerPrice      float64 `json:"tp,omitempty"` // trigger_price
	MarketProtection  float64 `json:"mp,omitempty"` // market_protection
	DisclosedQuantity int     `json:"dq,omitempty"` // disclosed_quantity
	FilledQuantity    int     `json:"fq,omitempty"` // filled_quantity
	DD                string  `json:"dd,omitempty"` // dd
	AMO               bool    `json:"am,omitempty"` // amo
}
