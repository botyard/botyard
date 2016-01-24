package config

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v2"
)

func Test_YamlConfig(t *testing.T) {
	var data = `name: boty
commands:
	- name: "hello"
	  command: "hello <world>"
	  url: "http://example.com/hello?world=<world>"
	  type: http
	  method: GET
irc:
	server: irc.org	
slack:
	token: TOKEN

`

	b := bytes.Replace([]byte(data), []byte("\t"), []byte("    "), -1)

	cfg := Config{}
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		t.Error(err)
	}

	if want, have := "hello <world>", cfg.Commands[0].Command; want != have {
		t.Errorf("want:%v,have:%v", want, have)
	}

}
