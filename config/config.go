package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Plugin defines the configuration for plugins to be loaded
type Plugin struct {
	Handler string `yaml:"handler"`
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
}

// Config defines the configuration items parsed from the config file
type Config struct {
	APIToken string   `yaml:"api_token"`
	BotName  string   `yaml:"bot_name"`
	Plugins  []Plugin `yaml:"plugins"`
}

// Parse reads the config file and stores the values in the struct
func (c *Config) Parse(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, c)
}
