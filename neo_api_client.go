package kotak_neo_api

import (
	"github.com/shikharvaish28/kotak-neo-api/api"
	"github.com/shikharvaish28/kotak-neo-api/websocket"
	"log"
	"context"
)

type KotakClient struct {
	config    api.Configuration
	websocket *websocket.HSWrapper // TODO: think about abstracting this away from the wrapper and building a simpler class instead.
}

// note: a broker interface should give you a channel for consumption and a client for placing orders.
func NewKotakClient(configuration api.Configuration) (*KotakClient, chan websocket.BrokerEvent) {
	ws, wsChannel := websocket.NewHSWrapper()
	return &KotakClient{
		config:    configuration,
		websocket: ws,
	}, wsChannel
}

// Subscribe handles the subscription to live feeds
func (c *KotakClient) Subscribe(ctx context.Context, instrumentTokens []string, isIndex, isDepth bool) {
	if c.config.EditToken == "" || c.config.EditSid == "" {
		log.Println("Please complete the Login Flow to Subscribe the Scrips")
		return
	}

	err := c.websocket.GetLiveFeed(instrumentTokens, isIndex, isDepth)
	if err != nil {
		log.Println("Failed to subscribe to live feeds:", err)
	}
}
