package ws

import (
	"github.com/centrifugal/centrifuge-go"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	client *centrifuge.Client
}

func NewWebSocketClient(server string) *Client {
	c := centrifuge.New(server, centrifuge.DefaultConfig())

	if err := c.Connect(); err != nil {
		log.Errorf("failed to connect to web socket: %s", err)
		return nil
	}

	return &Client{c}
}

func (e *Client) CreateSubscription(channel string) *centrifuge.Subscription {
	sub, err := e.client.NewSubscription(channel)
	if err != nil {
		log.Errorf("failed to create new subscription: %s", err)
		return nil
	}

	return sub
}

func (e Client) Subscribe(sub *centrifuge.Subscription) {
	if err := sub.Subscribe(); err != nil {
		log.Errorf("failed to subscribe: %s", err)
	}
}

func (e *Client) Close() {
	if err := e.client.Close(); err != nil {
		log.Errorf("failed to close ws connection: %s", err)
	}
}
