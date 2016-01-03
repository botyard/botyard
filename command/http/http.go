package http

import (
	"github.com/botyard/botyard/command/matcher"
	"github.com/botyard/botyard/command/parse"
	"github.com/botyard/botyard/message"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type HttpCommand struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Cmd    string `json:"cmd"`
	items  []*parse.Item
}

func New(method, url, cmd string) (*HttpCommand, error) {
	c := &HttpCommand{
		Method: method,
		Url:    url,
		Cmd:    cmd,
	}

	items, err := parse.Parse(c.Cmd)
	if err != nil {
		return nil, err
	}
	c.items = items
	return c, nil
}

func (c *HttpCommand) Match(in string) (req interface{}, ok bool) {
	m := matcher.New(c.Cmd, c.items)
	args, ok := m.MatchAndReturnArguments()
	return args, ok
}

func (c *HttpCommand) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		return nil, nil
	}
}

func (c *HttpCommand) Response(res interface{}) (msg *message.Message, err error) {

	return nil, nil
}
