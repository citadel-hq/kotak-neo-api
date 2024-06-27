# Kotak Neo Go SDK Development

- API version: 1.0.1
- Package version: 1.2.0

## Requirements.

Go 1.13+

## Installation & Usage
### go get

You can install the package using:

```sh
go get github.com/Kotak-Neo/kotak-neo-api/go
```

Then import the package:
```go
import neo "github.com/Kotak-Neo/kotak-neo-api/go"
```

### Go Modules

If you are using Go modules, simply import `github.com/Kotak-Neo/kotak-neo-api/go` in your code and run `go mod tidy`.

## Getting Started

Please follow the installation instructions and then run the following Go code:

```go
package main

import (
	"fmt"
	"github.com/Kotak-Neo/kotak-neo-api/go"
)

func main() {
	client := neo.NewClient("your_consumer_key", "your_consumer_secret", neo.UAT)

	// Initiate login by passing any of the combinations mobilenumber & password (or) pan & password (or) userid & password
	// Also this will generate the OTP to complete 2FA
	client.Login(neo.LoginOptions{
		MobileNumber: "+919999999999",
		Password:     "XXXX",
	})

	// Complete login and generate session token
	client.Session2FA("OTP")

	// Once 2FA has you can place the order by using below function
	orderID, err := client.PlaceOrder(neo.OrderOptions{
		ExchangeSegment: "NSE",
		Product:         "CNC",
		Price:           "100.50",
		OrderType:       "LIMIT",
		Quantity:        10,
		Validity:        "DAY",
		TradingSymbol:   "RELIANCE",
		TransactionType: "BUY",
	})
	if err != nil {
		fmt.Println("Error placing order:", err)
		return
	}
	fmt.Println("Order placed successfully, order ID:", orderID)
}
```

## Documentation for API Endpoints

All URIs are relative to *https://api.kotaksecurities.com/apim*

| Class             | Method                                                                        | HTTP request       | Description        |
|-------------------|-------------------------------------------------------------------------------|--------------------|--------------------|
| *LoginAPI*        | [**Login**](docs/Login.md#login)                                               | `POST /session/1.0/session/login` | Login NeoAPI       |
| *OrderAPI*        | [**PlaceOrder**](docs/Place_Order.md#place_order)                              | `POST /orders/3.0/orders` | Place Order        |
| *OrderAPI*        | [**ModifyOrder**](docs/Modify_Order.md#modify_order)                           | `PUT /orders/3.0/orders/{orderId}` | Modify Order       |
| *OrderAPI*        | [**CancelOrder**](docs/Cancel_Order.md#cancel_order)                           | `DELETE /orders/3.0/orders/{orderId}` | Cancel Order       |
| *OrderAPI*        | [**GetOrderReport**](docs/Order_report.md#order_report)                        | `GET /reports/3.0/orders` | Order Report       |
| *TradeAPI*        | [**GetTradeReport**](docs/Trade_report.md#trade_report)                        | `GET /reports/3.0/trades` | Trade Report       |
| *PortfolioAPI*    | [**GetPositions**](docs/Positions.md#positions)                                | `GET /portfolio/3.0/positions` | Positions          |
| *PortfolioAPI*    | [**GetHoldings**](docs/Holdings.md#holdings)                                   | `GET /portfolio/3.0/holdings` | Holdings           |
| *LimitsAPI*       | [**GetLimits**](docs/Limits.md#limits)                                         | `GET /user/3.0/limits` | Limits             |
| *MarginAPI*       | [**GetMarginRequired**](docs/Margin_Required.md#margin_required)               | `GET /orders/3.0/margins` | Margin Required    |
| *ScripMasterAPI*  | [**GetScripMaster**](docs/Scrip_Master.md#scrip_master)                        | `GET /scrip/3.0/scripmaster` | Scrip Master       |
| *ScripSearchAPI*  | [**SearchScrip**](docs/Scrip_Search.md#scrip_search)                           | `GET /scrip/3.0/search` | Scrip Search       |
| *QuotesAPI*       | [**GetQuotes**](docs/Quotes.md#quotes)                                         | `GET /market/3.0/quotes` | Quotes             |
| *WebSocketAPI*    | [**Subscribe**](docs/webSocket.md#subscribe)                                   | `WebSocket /ws/3.0` | Subscribe          |

