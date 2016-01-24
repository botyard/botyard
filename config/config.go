package config

type Config struct {
	Botname      string             `yaml:"name"`
	Commands     []CommandConfig    `yaml:"commands"`
	IRCGateway   IRCGatewayConfig   `yaml:"irc,omitempty"`
	SlackGateway SlackGatewayConfig `yaml:"slack,omitempty"`
}

type CommandConfig struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Command    string `yaml:"command"`
	HttpURL    string `yaml:"url,omitempty"`
	HttpMethod string `yaml:"method,omitempty"`
}

type IRCGatewayConfig struct {
	Server string `yaml:"server"`
	UseSSL bool   `yaml:"ssl,omitempty"`
}

type SlackGatewayConfig struct {
	Token string `yaml:"token"`
}
