package gateway

import (
	"crypto/tls"
	"log"

	"github.com/botyard/botyard/lib/config"
	"github.com/botyard/botyard/lib/message"

	irc "github.com/fluffle/goirc/client"
	"github.com/satori/go.uuid"
)

const (
	GATEWAY_IRC = "irc"
)

type IRCGateway struct {
	id         string
	cfg        config.IRCGateway
	msgChannel chan *message.Message
	channels   []string //#botyard...
	ircConn    *irc.Conn
}

func NewIRCGateway(cfg config.IRCGateway, botname string) *IRCGateway {

	//irc
	irccfg := irc.NewConfig(botname)
	if cfg.UseSSL == true {
		irccfg.SSL = true
		irccfg.SSLConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	irccfg.Server = cfg.Server
	irccfg.NewNick = func(n string) string { return n + "^" } //TODO:

	c := irc.Client(irccfg)

	gw := &IRCGateway{
		id:      uuid.NewV4().String(),
		cfg:     cfg,
		ircConn: c,
	}
	return gw
}

func (gw *IRCGateway) Open(c chan *message.Message) error {
	gw.msgChannel = c

	if err := gw.ircConn.Connect(); err != nil {
		log.Printf("Connection error: %s\n", err.Error())
		return err
	}

	gw.ircConn.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			for _, c := range gw.cfg.Channels {
				conn.Join(c)
				log.Printf("Joining %v\n", c)
			}
		})

	gw.ircConn.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			text := line.Text()
			channel := line.Target()
			m := message.FromGateway(
				gw.id,
				channel,
				text,
			)
			c <- m
		})

	log.Println("Open IRC.")

	return nil
}

func (gw *IRCGateway) SendMessage(m *message.Message) error {
	gw.ircConn.Privmsg(m.Address.Channel, m.Body)
	return nil //TODO:
}

func (gw *IRCGateway) ID() string {
	return gw.id
}

func (gw *IRCGateway) Name() string {
	return gw.cfg.Name
}
