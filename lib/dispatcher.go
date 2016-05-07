package lib

import (
	"github.com/botyard/botyard/lib/log"
	"github.com/botyard/botyard/lib/message"

	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

type Dispatcher struct {
	ctx        context.Context
	loader     *Loader
	msgChannel chan *message.Message
	logger     kitlog.Logger
}

func NewDispatcher(ctx context.Context,
	msgChannel chan *message.Message, loader *Loader) *Dispatcher {
	d := &Dispatcher{
		ctx:        ctx,
		loader:     loader,
		msgChannel: msgChannel,
		logger:     kitlog.NewContext(log.Logger).With("m", "dispatcher"),
	}

	go d.dispatch()
	return d
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case m := <-d.msgChannel:
			for _, c := range d.loader.Commands {
				req, ok := c.Match(m.Body)
				if ok {
					resp, err := c.Endpoint()(d.ctx, req)
					if err != nil {
						d.logger.Log("err", err)
						continue
					}
					reply, err := c.Response(resp)
					if err != nil {
						d.logger.Log("err", err)
						continue
					}
					reply.Address = m.Address //TODO:

					gw, ok := d.loader.Gateways[m.Address.GatewayID]
					if !ok {
						continue
					}
					gw.SendMessage(reply)
				}
			}
		}
	}
}
