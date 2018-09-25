package main

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

func Handler(args ...interface{}) {
	var (
		// c    config.Config
		api *slack.Client
		rtm *slack.RTM
		// info *slack.Info
		options map[string]interface{}
		ev      *slack.MessageEvent
	)

	for _, arg := range args {
		switch arg.(type) {
		// case config.Config:
		// 	c = arg.(config.Config)
		case *slack.Client:
			api = arg.(*slack.Client)
		case *slack.RTM:
			rtm = arg.(*slack.RTM)
		// case *slack.Info:
		// 	info = arg.(*slack.Info)
		case map[string]interface{}:
			options = arg.(map[string]interface{})
		case *slack.MessageEvent:
			ev = arg.(*slack.MessageEvent)
		}
	}

	msg := "Greetings!"
	if lookup, ok := options["lookup_user"]; ok {
		if lookup.(bool) {
			user, err := api.GetUserInfo(ev.User)
			if err != nil {
				log.Println("example_plugin: error getting user information:", err)
			}
			msg = fmt.Sprintf("Greetings, <@%s>!", user.ID)
		}
	}
	rtm.SendMessage(
		rtm.NewOutgoingMessage(msg, ev.Channel),
	)
}
