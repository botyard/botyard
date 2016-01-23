package main

import (
	"github.com/botyard/botyard/command"
	"github.com/botyard/botyard/command/http"
	"github.com/botyard/botyard/gateway"
	"github.com/botyard/botyard/message"
	"github.com/botyard/botyard/sync"
	"golang.org/x/net/context"

	"log"
)

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer func() { cancelFunc() }()

	var cmds map[string]command.Command
	{
		helloCmd, err := http.New("GET", "http://localhost:9001/hello", "hello")
		if err != nil {
			log.Fatal(err)
		}
		cmds["hello"] = helloCmd
	}

	var gateways map[string]gateway.Gateway
	{
		irc := gateway.NewIRCGateway()
		gateways[irc.ID()] = irc
	}

	var msgChannel chan *message.Message
	{
		msgChannel = make(chan *message.Message)
	}

	var dispatchers []*Dispatcher
	{
		for i := 0; i < 1; i++ {
			dispatcher := NewDispatcher(ctx, msgChannel, gateways, cmds)
			dispatchers = append(dispatchers, dispatcher)
			sync.WaitGroup.Add(1)
		}
	}

	sync.WaitGroup.Wait()
}
