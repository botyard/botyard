package gateway

import (
	"crypto/tls"
	"log"

	"github.com/botyard/botyard/message"
	irc "github.com/fluffle/goirc/client"
	"github.com/satori/go.uuid"
)

const (
	GATEWAY_IRC = "irc"
)

type IRCGateway struct {
	id         string
	msgChannel chan *message.Message
	channels   []string //#botyard...
	ircConn    *irc.Conn
}

func NewIRCGateway() *IRCGateway {
	//irc
	cfg := irc.NewConfig("boty")
	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	cfg.Server = "irc.ozinger.org:16667" //TODO
	cfg.NewNick = func(n string) string { return n + "^" }
	c := irc.Client(cfg)

	gw := &IRCGateway{
		id:      uuid.NewV4().String(),
		ircConn: c,
	}
	return gw
}

func (gw *IRCGateway) Open(c chan *message.Message) {
	gw.msgChannel = c

	log.Println("PRE connect")
	if err := gw.ircConn.Connect(); err != nil {
		log.Printf("Connection error: %s\n", err.Error())
		return
	}
	log.Printf("Post connect")

	gw.ircConn.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) { conn.Join("#botyard") })

	gw.ircConn.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			text := line.Text()
			m := message.FromGateway(
				gw.id,
				"#botyard",
				text,
			)
			c <- m
		})
	log.Println("Open..")
}

func (gw *IRCGateway) SendMessage(m *message.Message) error {
	gw.ircConn.Privmsg(m.Address.Channel, m.Body)
	return nil //TODO:
}

func (gw *IRCGateway) ID() string {
	return gw.id
}

//TODO

/*
func init() {
	gw := NewIRCGateway()
	Add(gw.ID, gw)
}
*/
