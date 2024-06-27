// Package api implements the order history API in Go.
package go

import (
	"encoding/json"
	"net/http"
)

// OrderHistoryService provides operations to fetch order history.
type OrderHistoryService struct {
	client *Client
}

// NewOrderHistoryService creates a new service for order history.
func NewOrderHistoryService(client *Client) *OrderHistoryService {
	return &OrderHistoryService{
		client: client,
	}
}

// FetchOrderHistory fetches the history of an order.
func (s *OrderHistoryService) FetchOrderHistory(orderID string) (*OrderHistoryResponse, error) {
	req, err := s.client.NewRequest("GET", "/order_history/"+orderID, nil)
	if err != nil {
		return nil, err
	}

	var resp OrderHistoryResponse
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// OrderHistoryResponse represents a response from an order history operation.
type OrderHistoryResponse struct {
	// Add order history response fields here
}
