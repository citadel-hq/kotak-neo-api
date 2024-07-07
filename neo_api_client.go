package kotak_neo_api

import (
	"github.com/shikharvaish28/kotak-neo-api/api"
	"github.com/shikharvaish28/kotak-neo-api/websocket"
	"log"
	"context"
	"errors"
)

type KotakClient struct {
	config    api.Configuration
	websocket *websocket.HSWrapper // TODO: think about abstracting this away from the wrapper and building a simpler class instead.
	loginAPI  *api.LoginAPI
	orderAPI  *api.OrderService
}

// note: a broker interface should give you a channel for consumption and a client for placing orders.
func NewKotakClient(configuration api.Configuration) (*KotakClient, chan websocket.BrokerEvent) {
	ws, wsChannel := websocket.NewHSWrapper()
	apiClient := api.NewAPIClient(configuration)
	loginAPI := api.NewLoginAPI(apiClient)
	orderAPI := api.NewOrderService(apiClient)
	return &KotakClient{
		config:    configuration,
		websocket: ws,
		loginAPI:  loginAPI,
		orderAPI:  orderAPI,
	}, wsChannel
}

// Subscribe handles the subscription to live feeds
func (c *KotakClient) Subscribe(ctx context.Context, instrumentTokens []map[string]string, isIndex, isDepth bool) {
	if c.config.EditToken == "" || c.config.EditSid == "" {
		log.Println("Please complete the Login Flow to Subscribe the Scrips")
		return
	}

	err := c.websocket.GetLiveFeed(instrumentTokens, isIndex, isDepth)
	if err != nil {
		log.Println("Failed to subscribe to live feeds:", err)
	}
}

func (c *KotakClient) Login(password, mobileNumber, userid, pan, MPin string) (map[string]interface{}, error) {
	if mobileNumber == "" && userid == "" && pan == "" {
		errorResponse := map[string]interface{}{
			"error": []map[string]string{
				{"code": "10300", "message": "Validation Errors! Any of Mobile Number, User Id and Pan has to pass as part of login"},
			},
		}
		return errorResponse, errors.New("validation error")
	}

	viewToken, err := c.loginAPI.GenerateViewToken(password, mobileNumber, userid, pan, MPin)
	if err != nil {
		return nil, err
	}

	_, err = c.loginAPI.GenerateOTP()
	if err != nil {
		return map[string]interface{}{
			"error": []map[string]string{
				{"code": "10522", "message": "Issues while generating OTP! Try to login again."},
			},
		}, errors.New("issues while generating OTP")
	}

	return viewToken, nil
}

func (c *KotakClient) Session2FA(otp string) (map[string]interface{}, error) {
	editToken, err := c.loginAPI.Login2FA(otp)
	if err != nil {
		return nil, err
	}
	return editToken, nil
}

func (c *KotakClient) PlaceOrder(req api.OrderRequest) (map[string]interface{}, error) {
	if c.config.EditToken != "" && c.config.EditSid != "" {
		// TODO: insert the validation here
		//exchangeSegment := api.ExchangeSegment[req.ExchangeSegment]
		//product := api.Product[req.Product]
		//orderType := api.OrderType[req.OrderType]
		return c.orderAPI.PlaceOrder(req)
	}
	return map[string]interface{}{
		"error": []map[string]string{
			{"code": "900", "message": "Complete the 2fa process before accessing this application"},
		},
	}, errors.New("please complete the 2FA process before placing orders")

}
