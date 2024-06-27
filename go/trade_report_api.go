// Package go provides the Go implementation for trade report API operations.
package go

import (
	"encoding/json"
	"net/http"
)

// TradeReportAPI struct defines the client for the trade report API operations.
type TradeReportAPI struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewTradeReportAPI creates a new trade report API client with the specified base URL.
func NewTradeReportAPI(baseURL string) *TradeReportAPI {
	return &TradeReportAPI{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// GetTradeReport fetches the trade report for all trades or a specific trade by order ID.
func (api *TradeReportAPI) GetTradeReport(orderID string) ([]byte, error) {
	url := api.BaseURL + "/trade_report"
	if orderID != "" {
		url += "?order_id=" + orderID
	}

	resp, err := api.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
