package main

import (
	"flag"
	"fmt"
	"log"
	"plugin"
	"regexp"
	"strings"

	"github.com/mattjw79/plugbot/config"
	"github.com/nlopes/slack"
)

// Plugin holds the loaded plugin handler
type Plugin func(...interface{})

// PlugBot defines the main bot structure
type PlugBot struct {
	ConfigFile string
	Config     config.Config
	API        *slack.Client
	RTM        *slack.RTM
	Info       *slack.Info
	Plugins    map[string]Plugin
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

// Init loads prereqs and starts the Slack connection
func (p *PlugBot) Init() error {
	p.Plugins = make(map[string]Plugin)
	p.ParseFlags()
	if err := p.ParseConfig(); err != nil {
		return err
	}
	p.LoadPlugins()
	p.API = slack.New(p.Config.APIToken)
	p.RTM = p.API.NewRTM()
	go p.RTM.ManageConnection()
	return nil
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

// LoadPlugins loads the plugins defined
func (p *PlugBot) LoadPlugins() {
	for _, configPlugin := range p.Config.Plugins {
		plug, err := plugin.Open(configPlugin.Path)
		if err != nil {
			log.Println("error loading plugin:", err)
			continue
		}

		handler, err := plug.Lookup(configPlugin.Handler)
		if err != nil {
			log.Println("error loading plugin handler:", err)
			continue
		}

		log.Printf("loaded plugin '%s'\n", configPlugin.Name)
		p.Plugins[configPlugin.Name] = handler.(func(...interface{}))
	}
}

func main() {
	var plugbot PlugBot
	if err := plugbot.Init(); err != nil {
		log.Fatal("error initializing:", err)
	}

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
			for _, eventConfig := range plugbot.Config.Events {
				regex := eventConfig.Regex
				if strings.Contains(regex, "<@self>") {
					regex = strings.Replace(eventConfig.Regex, "<@self>", fmt.Sprintf("<@%s>", plugbot.Info.User.ID), -1)
				}
				matched, err := regexp.MatchString(regex, ev.Text)
				if err != nil {
					log.Printf("error matching regex:\n  `%s`\n  '%s'", regex, ev.Text)
				}
				if matched {
					if !eventConfig.Recipient || (eventConfig.Recipient && plugbot.IsRecipient(ev)) {
						if eventConfig.Response != "" {
							plugbot.RTM.SendMessage(
								plugbot.RTM.NewOutgoingMessage(eventConfig.Response, ev.Channel),
							)
						}
					}
				}
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
