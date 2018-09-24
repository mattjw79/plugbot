package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config defines the configuration items parsed from the config file
type Config struct {
	APIToken string `yaml:"api_token"`
	BotName  string `yaml:"bot_name"`
}

// Parse reads the config file and stores the values in the struct
func (c *Config) Parse(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, c)
}
