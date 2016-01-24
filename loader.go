package main

import (
	"github.com/botyard/botyard/command"
	httpcmd "github.com/botyard/botyard/command/http"
	"github.com/botyard/botyard/config"
	"github.com/botyard/botyard/gateway"
	"github.com/botyard/botyard/message"

	"gopkg.in/yaml.v2"

	"bytes"
	"fmt"
	"log"
)

type Loader struct {
	Config   config.Config
	Commands map[string]command.Command
	Gateways map[string]gateway.Gateway
}

func NewLoader(data []byte) (*Loader, error) {
	l := &Loader{
		Commands: make(map[string]command.Command),
		Gateways: make(map[string]gateway.Gateway),
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

			cmd, err := httpcmd.New(cmdcfg.HttpMethod, cmdcfg.HttpURL, cmdcfg.Command)
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
		log.Printf("Load irc gateway id:%v name:%v", irc.ID(), irccfg.Name)
	}
	for _, slackcfg := range l.Config.SlackGateway {
		//TODO: valudate config
		slack := gateway.NewSlackGateway(slackcfg)
		l.Gateways[slack.ID()] = slack
		log.Printf("Load slack gateway id:%v name:%v", slack.ID(), slackcfg.Name)
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
