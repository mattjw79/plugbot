package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/nlopes/slack"
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

// PlugBot defines the main bot structure
type PlugBot struct {
	ConfigFile string
	Config     Config
	API        *slack.Client
	RTM        *slack.RTM
	Info       *slack.Info
}

// ParseFlags parses the command line flags
func (p *PlugBot) ParseFlags() {
	flag.StringVar(&p.ConfigFile, "config", "config.yml", "Path to the config file")
	flag.Parse()
}

// ParseConfig runs the config file parsing
func (p *PlugBot) ParseConfig() error {
	return p.Config.Parse(p.ConfigFile)
}

// IsRecipient determines if an incoming message event is intended for the bot
func (p *PlugBot) IsRecipient(msg *slack.MessageEvent) bool {
	isPrivate := false
	_, cErr := p.API.GetChannelInfo(msg.Channel)
	_, gErr := p.API.GetGroupInfo(msg.Channel)
	if cErr != nil && gErr != nil {
		isPrivate = true
	}

	if strings.Contains(msg.Text, fmt.Sprintf("<@%s>", p.Info.User.ID)) || isPrivate {
		return true
	}
	return false
}

func main() {
	var plugbot PlugBot
	plugbot.ParseFlags()
	if err := plugbot.ParseConfig(); err != nil {
		log.Fatal("error parsing config:", err)
	}
	plugbot.API = slack.New(plugbot.Config.APIToken)
	plugbot.RTM = plugbot.API.NewRTM()
	go plugbot.RTM.ManageConnection()

	for msg := range plugbot.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			plugbot.Info = plugbot.RTM.GetInfo()
			log.Printf(
				"connected to Slack team '%s' (%s) as '%s' (%s)\n",
				plugbot.Info.Team.Name,
				plugbot.Info.Team.ID,
				plugbot.Info.User.Name,
				plugbot.Info.User.ID,
			)
		case *slack.MessageEvent:
			if plugbot.IsRecipient(ev) {
				plugbot.RTM.SendMessage(
					plugbot.RTM.NewOutgoingMessage(fmt.Sprintf("Message received: %v", ev.Text), ev.Channel),
				)
			}
		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			log.Fatal("invalid credentials")
			return
		default:
		}
	}
}