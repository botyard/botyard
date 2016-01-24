package command

import (
	"github.com/botyard/botyard/message"

	"github.com/go-kit/kit/endpoint"
)

const (
	HTTP = "http"
)

type Command interface {
	Match(in string) (req interface{}, ok bool)
	Endpoint() endpoint.Endpoint
	Response(interface{}) (*message.Message, error)
}
