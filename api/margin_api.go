// Package api provides the Go implementation for margin API operations.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MarginAPI struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

func NewMarginAPI(apiClient *APIClient) *MarginAPI {
	return &MarginAPI{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

func (api *MarginAPI) GetMargin(request MarginRequest) (map[string]interface{}, error) {
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

	url, err := api.ApiClient.Config.getUrlDetails("margin")
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

// MarginRequest represents the request structure for margin calculation.
type MarginRequest struct {
	ExchangeSegment  string  `json:"exSeg,omitempty"`
	Price            float64 `json:"price,omitempty"`
	OrderType        string  `json:"prcTp,omitempty"`
	Product          string  `json:"prod,omitempty"`
	Quantity         int     `json:"qty,omitempty"`
	InstrumentToken  string  `json:"tok,omitempty"`
	TransactionType  string  `json:"trnsTp,omitempty"`
	TriggerPrice     float64 `json:"trgPrc,omitempty"`
	BrokerName       string  `json:"brkName,omitempty"`
	BranchId         string  `json:"brnchId,omitempty"`
	StopLossType     string  `json:"slAbsOrTks,omitempty"`
	StopLossValue    string  `json:"slVal,omitempty"`
	SquareOffType    string  `json:"sqrOffAbsOrTks,omitempty"`
	SquareOffValue   float64 `json:"sqrOffVal,omitempty"`
	TrailingStopLoss float64 `json:"trailSL,omitempty"`
	TrailingSLValue  float64 `json:"tSLTks,omitempty"`
}
