package main

import (
	"github.com/botyard/botyard/command"
	"github.com/botyard/botyard/gateway"
	"github.com/botyard/botyard/message"
	"golang.org/x/net/context"

	"log"
)

type Dispatcher struct {
	ctx        context.Context
	Commands   map[string]command.Command
	Gateways   map[string]gateway.Gateway
	msgChannel chan *message.Message
}

func NewDispatcher(ctx context.Context,
	msgChannel chan *message.Message,
	gws map[string]gateway.Gateway,
	cmds map[string]command.Command) *Dispatcher {
	d := &Dispatcher{
		ctx:        ctx,
		Commands:   cmds,
		Gateways:   gws,
		msgChannel: msgChannel,
	}

	go d.dispatch()
	return d
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case m := <-d.msgChannel:
			log.Println(m)
			for _, c := range d.Commands {
				log.Println("body:", m.Body)
				req, ok := c.Match(m.Body)
				if ok {
					resp, err := c.Endpoint()(d.ctx, req)
					if err != nil {
						log.Println(err)
						continue
					}
					reply, err := c.Response(resp)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println("reply:", reply)
					reply.Address = m.Address //TODO:

					log.Println("reply.addr:", reply.Address)

					gw, ok := d.Gateways[m.Address.GatewayID]
					if !ok {
						continue
					}
					log.Println("reply gw:", m.Address.GatewayID)
					gw.SendMessage(reply)

				}
			}
		}
	}
}
