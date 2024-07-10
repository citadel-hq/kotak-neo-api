// Package api provides the Go implementation for Scrip Search API operations.
package api

// ScripSearchService handles communication with the scrip search related methods of the API.
type ScripSearchService struct {
	ApiClient  *APIClient
	RestClient *RESTClientObject
}

// NewScripSearchService creates a new ScripSearchService.
func NewScripSearchService(apiClient *APIClient) *ScripSearchService {
	return &ScripSearchService{
		ApiClient:  apiClient,
		RestClient: apiClient.RestClient,
	}
}

// ScripSearchRequest represents a request to the scrip search API.
type ScripSearchRequest struct {
	Keyword          string  `json:"keyword"`
	Symbol           string  `json:"symbol"`
	ExchangeSegment  string  `json:"exchangeSegment"`
	Expiry           string  `json:"expiry"`
	OptionType       string  `json:"optionType"`
	StrikePrice      float64 `json:"strikePrice"`
	Ignore50Multiple bool    `json:"ignore50Multiple"`
}

// TODO: need to complete this, slightly complicated.
//func (api *ScripSearchService) ScripSearch(request *ScripSearchRequest) (map[string]interface{}, error) {
//	headerParams := map[string]string{
//		"Authorization": fmt.Sprintf("Bearer %s", api.ApiClient.Config.BearerToken),
//	}
//	url, err := api.ApiClient.Config.getUrlDetails("scrip_master")
//	if err != nil {
//		return nil, err
//	}
//	resp, err := api.RestClient.Request(http.MethodGet, url, nil, headerParams, nil)
//	if err != nil {
//		return nil, err
//	}
//	var scripReport map[string]interface{}
//	_ = json.NewDecoder(resp.Body).Decode(&scripReport)
//	data := scripReport["data"].(map[string]interface{})
//
//	var csvFileURL string
//	if request.ExchangeSegment != "" {
//		for _, file := range data["filesPaths"].([]interface{}) {
//			if strings.Contains(strings.ToLower(file.(string)), strings.ToLower(request.ExchangeSegment)) {
//				csvFileURL = file.(string)
//				break
//			}
//		}
//	}
//
//	if csvFileURL == "" {
//		return map[string]interface{}{
//			"error": []map[string]string{
//				{"code": "10300", "message": "No matching CSV file found for the exchange segment."},
//			},
//		}, nil
//	}
//
//	resp, err = http.Get(csvFileURL)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	reader := csv.NewReader(resp.Body)
//	records, err := gocsv.UnmarshalCSVWithoutHeaders(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	// Filtering logic
//	//filteredRecords := filterRecords(records, symbol, exchangeSegment, expiry, optionType, strikePrice)
//
//	if len(filteredRecords) == 0 {
//		return map[string]interface{}{
//			"message": "No data found with the given search information. Please try with other combinations.",
//		}, nil
//	}
//
//	filteredRecordsJSON, err := json.Marshal(filteredRecords)
//	if err != nil {
//		return nil, err
//	}
//
//	var result interface{}
//	if err := json.Unmarshal(filteredRecordsJSON, &result); err != nil {
//		return nil, err
//	}
//
//	return nil, nil
//}
//
//// filterRecords filters the CSV records based on the given criteria
//func filterRecords(records []*Record, symbol, exchangeSegment, expiry, optionType, strikePrice string) []*Record {
//	var filteredRecords []*Record
//	for _, record := range records {
//		if strings.Contains(strings.ToLower(record.PSymbolName), strings.ToLower(symbol)) {
//			if optionType != "" {
//				if !strings.Contains(strings.ToLower(record.POptionType), strings.ToLower(optionType)) {
//					continue
//				}
//			}
//			if expiry != "" {
//				expiryDates := strings.Split(expiry, "-")
//				expiryDate, _ := time.Parse("02Jan2006", record.PExpiryDate)
//				if len(expiryDates) == 2 {
//					startDate, _ := time.Parse("02Jan2006", expiryDates[0])
//					endDate, _ := time.Parse("02Jan2006", expiryDates[1])
//					if expiryDate.Before(startDate) || expiryDate.After(endDate) {
//						continue
//					}
//				} else {
//					singleDate, _ := time.Parse("02Jan2006", expiryDates[0])
//					if !expiryDate.Equal(singleDate) {
//						continue
//					}
//				}
//			}
//			if strikePrice != "" {
//				// Implement strike price filtering logic
//			}
//			filteredRecords = append(filteredRecords, record)
//		}
//	}
//	return filteredRecords
//}
//
//// Record represents a CSV record
//type Record struct {
//	PSymbolName string `csv:"pSymbolName"`
//	POptionType string `csv:"pOptionType"`
//	PExpiryDate string `csv:"pExpiryDate"`
//}
