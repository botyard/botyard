package lib

import (
	"github.com/botyard/botyard/lib/command"
	httpcmd "github.com/botyard/botyard/lib/command/http"
	"github.com/botyard/botyard/lib/config"
	"github.com/botyard/botyard/lib/gateway"
	"github.com/botyard/botyard/lib/log"
	"github.com/botyard/botyard/lib/message"

	kitlog "github.com/go-kit/kit/log"
	"gopkg.in/yaml.v2"

	"bytes"
	"fmt"
)

type Loader struct {
	Config   config.Config
	Commands map[string]command.Command
	Gateways map[string]gateway.Gateway
	logger   kitlog.Logger
}

func NewLoader(data []byte) (*Loader, error) {
	l := &Loader{
		Commands: make(map[string]command.Command),
		Gateways: make(map[string]gateway.Gateway),
		logger:   kitlog.NewContext(log.Logger).With("m", "Loader"),
	}

	cfg, err := l.loadConfigFile(data)
	if err != nil {
		return nil, err
	}
	l.Config = cfg

	if err := l.loadModules(); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *Loader) loadModules() error {
	if err := l.loadCommands(); err != nil {
		return err
	}
	if err := l.loadGateways(); err != nil {
		return err
	}

	return nil
}

func (l *Loader) loadCommands() error {
	for _, cmdcfg := range l.Config.Commands {
		switch cmdcfg.Type {
		case command.HTTP:
			//TODO: valudate config e.g. err := config.Validate()
			if cmdcfg.Name == "" {
				return fmt.Errorf("The command config hasn't name")
			}

			var cmdStr string
			matchWords := false
			if cmdcfg.Command != "" {
				cmdStr = fmt.Sprintf("%s %s", l.Config.Botname, cmdcfg.Command)
			} else if cmdcfg.Words != "" {
				cmdStr = cmdcfg.Words
				matchWords = true
			} else {
				return fmt.Errorf("The command needs to have one of command or words")
			}

			cmd, err := httpcmd.New(cmdcfg.HttpMethod, cmdcfg.HttpURL, cmdStr, matchWords)
			if err != nil {
				return err
			}

			l.Commands[cmdcfg.Name] = cmd
		}
	}

	return nil
}

func (l *Loader) loadGateways() error {
	for _, irccfg := range l.Config.IRCGateway {
		//TODO: valudate config
		irc := gateway.NewIRCGateway(irccfg, l.Config.Botname)
		l.Gateways[irc.ID()] = irc
		l.logger.Log("load", "irc", "id", irc.ID(), "name", irccfg.Name)
	}
	for _, slackcfg := range l.Config.SlackGateway {
		//TODO: valudate config
		slack := gateway.NewSlackGateway(slackcfg)
		l.Gateways[slack.ID()] = slack
		l.logger.Log("load", "slack", "id", slack.ID(), "name", slackcfg.Name)
	}
	return nil
}

func (l *Loader) loadConfigFile(data []byte) (config.Config, error) {
	cfg := config.Config{}

	data = bytes.Replace([]byte(data), []byte("\t"), []byte("    "), -1)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

//TODO: Is it right place for openning gateway?...
func (l *Loader) OpenGateways(msgCh chan *message.Message) error {
	for _, gw := range l.Gateways {
		if err := gw.Open(msgCh); err != nil {
			return err
		}
	}
	return nil
}
