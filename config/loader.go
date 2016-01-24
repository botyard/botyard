package config

import (
	"gopkg.in/yaml.v2"

	"bytes"
	"io/ioutil"
)

func LoadConfigFile(filename string) (*Config, error) {
	config := Config{}

	yamlContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	yamlContent = bytes.Replace([]byte(yamlContent), []byte("\t"), []byte("    "), -1)

	err = yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
