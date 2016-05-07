package gateway

import (
	"github.com/botyard/botyard/lib/message"
)

type Gateway interface {
	ID() string
	Name() string
	Open(chan *message.Message) error
	SendMessage(*message.Message) error
	//Close() //TODO:
}
