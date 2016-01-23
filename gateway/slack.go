package gateway

import (
	"github.com/botyard/botyard/message"

	"github.com/bobbytables/slacker"

	"fmt"
	"log"
)

const TOKEN = "xoxb-19244391940-ZrKpxFdwo1F0a4P8uQpO5U2e"

type SlackGateway struct {
	id     string
	client *slacker.APIClient
	broker *slacker.RTMBroker
}

func NewSlackGateway() *SlackGateway {
	client := slacker.NewAPIClient(TOKEN, "")
	gw := &SlackGateway{
		id:     "slack",
		client: client,
	}
	return gw
}

func (gw *SlackGateway) ID() string {
	return gw.id
}

func (gw *SlackGateway) Open(c chan *message.Message) {
	rtmStart, err := gw.client.RTMStart()
	if err != nil {
		log.Println(err)
		return
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
