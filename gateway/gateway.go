package gateway

import (
	"github.com/botyard/botyard/message"
)

//var Gateways = map[string]Gateway{}

type Gateway interface {
	ID() string
	Open(chan *message.Message)
	SendMessage(*message.Message) error
	//Close()
}

/*
func Add(id string, gw Gateway) {
	Gateways[id] = gw
}
*/
