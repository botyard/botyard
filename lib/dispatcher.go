package lib

import (
	"github.com/botyard/botyard/lib/message"

	"golang.org/x/net/context"

	"log"
)

type Dispatcher struct {
	ctx        context.Context
	loader     *Loader
	msgChannel chan *message.Message
}

func NewDispatcher(ctx context.Context,
	msgChannel chan *message.Message, loader *Loader) *Dispatcher {
	d := &Dispatcher{
		ctx:        ctx,
		loader:     loader,
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
			for _, c := range d.loader.Commands {
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

					gw, ok := d.loader.Gateways[m.Address.GatewayID]
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
