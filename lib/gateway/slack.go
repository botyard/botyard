package gateway

import (
	"github.com/botyard/botyard/lib/config"
	"github.com/botyard/botyard/lib/log"
	"github.com/botyard/botyard/lib/message"

	"github.com/bobbytables/slacker"
	kitlog "github.com/go-kit/kit/log"
)

const (
	SlackEventTypeMessage = "message"
)

type SlackGateway struct {
	id     string
	cfg    config.SlackGateway
	client *slacker.APIClient
	broker *slacker.RTMBroker
	logger kitlog.Logger
}

func NewSlackGateway(cfg config.SlackGateway) *SlackGateway {
	client := slacker.NewAPIClient(cfg.Token, "")
	gw := &SlackGateway{
		id:     "slack",
		cfg:    cfg,
		client: client,
		logger: kitlog.NewContext(log.Logger).With("m", "Gateway-Slack"),
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

			if event.Type == SlackEventTypeMessage {
				msg, err := event.Message()
				if err != nil {
					gw.logger.Log("err", err, "event", event)
					panic(err)
				}

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

	err := gw.broker.Publish(rtmMsg)
	return err

}
