package go

import (
	"encoding/json"
	"net/http"
)

// OrderReportAPI handles the order report functionality
type OrderReportAPI struct {
	Client *APIClient
}

// NewOrderReportAPI creates a new OrderReportAPI instance
func NewOrderReportAPI(client *APIClient) *OrderReportAPI {
	return &OrderReportAPI{
		Client: client,
	}
}

// GetOrderReport fetches the order report for the given order ID
func (api *OrderReportAPI) GetOrderReport(orderID string) (*OrderReport, error) {
	resp, err := api.Client.Request("GET", "/order/report/"+orderID, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var report OrderReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, err
	}

	return &report, nil
}

// OrderReport represents the structure of an order report
type OrderReport struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	Description string `json:"description"`
	// Add other fields as necessary
}
