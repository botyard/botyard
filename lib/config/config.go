package config

type Config struct {
	Botname      string         `yaml:"name"`
	Commands     []Command      `yaml:"commands"`
	IRCGateway   []IRCGateway   `yaml:"irc,omitempty"`
	SlackGateway []SlackGateway `yaml:"slack,omitempty"`
}

type Command struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Command    string `yaml:"command"`
	Words      string `yaml:"words"`
	HttpURL    string `yaml:"url,omitempty"`
	HttpMethod string `yaml:"method,omitempty"`
	Desc       string `yaml:"desc,omitempty"`
}

type IRCGateway struct {
	Name     string   `yaml:"name"`
	Server   string   `yaml:"server"`
	UseSSL   bool     `yaml:"ssl,omitempty"`
	Channels []string `yaml:"channels"`
}

type SlackGateway struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}
