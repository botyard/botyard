package message

import (
	"fmt"
)

type Message struct {
	Address *Address `json:"address",omitempty`
	Body    string   `json:"body"`
}

type Address struct {
	GatewayID string `json:"gateway_id"`
	Channel   string `json:"channel"`
}

func FromGateway(gatewayID, channel, body string) *Message {
	m := &Message{
		Address: &Address{
			GatewayID: gatewayID,
			Channel:   channel,
		},
		Body: body,
	}
	return m
}

func (a Address) String() string {
	return fmt.Sprintf("%s", a.GatewayID)
}
