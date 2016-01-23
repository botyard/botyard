package main

import (
	"github.com/botyard/botyard/command"
	httpcmd "github.com/botyard/botyard/command/http"
	"github.com/botyard/botyard/gateway"
	"github.com/botyard/botyard/message"
	"github.com/botyard/botyard/sync"
	"golang.org/x/net/context"

	"log"
	"net/http"
)

func main() {
	log.Println("Start.")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer func() { cancelFunc() }()

	var msgChannel chan *message.Message
	{
		msgChannel = make(chan *message.Message)
	}

	var cmds map[string]command.Command
	{
		cmds = make(map[string]command.Command)
		helloCmd, err := httpcmd.New("GET", "http://localhost:9001/hello?world=<world>", "hello <world>")
		if err != nil {
			log.Fatal(err)
		}
		cmds["hello"] = helloCmd
	}

	var gateways map[string]gateway.Gateway
	{
		gateways = make(map[string]gateway.Gateway)
		//irc := gateway.NewIRCGateway()
		//irc.Open(msgChannel)
		slack := gateway.NewSlackGateway()
		slack.Open(msgChannel)
		gateways[slack.ID()] = slack
		//gateways[irc.ID()] = irc
	}

	var dispatchers []*Dispatcher
	{
		for i := 0; i < 1; i++ {
			dispatcher := NewDispatcher(ctx, msgChannel, gateways, cmds)
			dispatchers = append(dispatchers, dispatcher)
			sync.WaitGroup.Add(1)
		}
	}

	go func() {
		http.HandleFunc("/hello", helloHandler)
		http.ListenAndServe(":9001", nil)
	}()
	sync.WaitGroup.Wait()
	log.Println("End.")
}
