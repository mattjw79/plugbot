package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Events defines the events and how to react
type Events struct {
	Options   map[string]interface{} `yaml:"options,omitempty"`
	Plugin    string                 `yaml:"plugin,omitempty"`
	Recipient bool                   `yaml:"recipient"`
	Regex     string                 `yaml:"regex"`
	Response  string                 `yaml:"response,omitempty"`
	Type      string                 `yaml:"type"`
}

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
	Events   []Events `yaml:"events"`
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
