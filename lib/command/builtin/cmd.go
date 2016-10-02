package builtin

import (
	"github.com/botyard/botyard/lib/command"
	"github.com/botyard/botyard/lib/command/matcher"
	"github.com/botyard/botyard/lib/command/parse"
	"github.com/botyard/botyard/lib/log"
	"github.com/botyard/botyard/lib/message"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"fmt"
)

type CmdFunc func(args []*command.Argument) (string, error)

type BuiltinCommand struct {
	Cmd        string
	CmdFunc    CmdFunc
	matchWords bool
	items      []*parse.Item
	logger     kitlog.Logger
}

func New(cmd string, cmdFunc CmdFunc, matchWords bool) (*BuiltinCommand, error) {
	c := &BuiltinCommand{
		Cmd:        cmd,
		CmdFunc:    cmdFunc,
		matchWords: matchWords,
		logger:     kitlog.NewContext(log.Logger).With("m", "BuiltinCommand"),
	}

	items, err := parse.Parse(c.Cmd)
	if err != nil {
		return nil, err
	}
	c.items = items
	return c, nil
}

func (c *BuiltinCommand) Match(in string) (req interface{}, ok bool) {
	cfg := &matcher.Config{}
	if c.matchWords == false {
		cfg.FitForward = true
		cfg.FitBackward = true
	}
	m := matcher.New(in, c.items, cfg)
	ok, _ = m.Match()
	req = m.Arguments()
	return
}

func (c *BuiltinCommand) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		args, ok := req.([]*command.Argument)
		if !ok {
			return nil, fmt.Errorf("TypeError: The request is not argument")
		}
		return c.CmdFunc(args)
	}
}

func (c *BuiltinCommand) Response(res interface{}) (*message.Message, error) {
	response, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("res is not http.Response")
	}

	msg := &message.Message{}
	msg.Body = response

	return msg, nil
}
