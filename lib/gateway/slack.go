package gateway

import (
	"github.com/botyard/botyard/lib/config"
	"github.com/botyard/botyard/lib/message"

	"github.com/bobbytables/slacker"

	"fmt"
	"log"
)

type SlackGateway struct {
	id     string
	cfg    config.SlackGateway
	client *slacker.APIClient
	broker *slacker.RTMBroker
}

func NewSlackGateway(cfg config.SlackGateway) *SlackGateway {
	client := slacker.NewAPIClient(cfg.Token, "")
	gw := &SlackGateway{
		id:     "slack",
		cfg:    cfg,
		client: client,
	}
	return gw
}

func (gw *SlackGateway) ID() string {
	return gw.id
}

func (gw *SlackGateway) Name() string {
	return gw.cfg.Name
}

func (gw *SlackGateway) Open(c chan *message.Message) error {
	rtmStart, err := gw.client.RTMStart()
	if err != nil {
		return err
	}
	gw.broker = slacker.NewRTMBroker(rtmStart)
	gw.broker.Connect()

	go func() {
		for {
			event := <-gw.broker.Events()
			fmt.Println(event.Type)

			if event.Type == "message" {
				msg, err := event.Message()
				if err != nil {
					panic(err)
				}

				fmt.Println(msg.Text)
				m := message.FromGateway(
					gw.id,
					msg.Channel,
					msg.Text,
				)
				c <- m

			}

		}
	}()

	return nil

}

func (gw *SlackGateway) SendMessage(m *message.Message) error {
	rtmMsg := &slacker.RTMMessage{}
	rtmMsg.Type = "message"
	rtmMsg.Text = m.Body
	rtmMsg.Channel = m.Address.Channel
	rtmMsg.User = "boty"
	//rtmMsg.Ts =

	log.Println("rtmMsg:", rtmMsg)
	err := gw.broker.Publish(rtmMsg)
	log.Println("SendMessage:", err)
	return err

}
