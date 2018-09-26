# PlugBot

PlugBot is designed to be an extensible framework for creating a Slack chat-bot.

## Building PlugBot

```
go get github.com/mattjw79/plugbot
go build plugbot.go
```

## Building plugins

```
go build -buildmode=plugin <plugin>.go
```

## Configuration

Configuration is done with a YAML file (default config.yml).

```
---

api_token: xoxb-00000000000-000000000000-000000000000000000000000
bot_name: PlugBot
plugins:
  - handler: Handler
    name: Example
    path: plugins/example_plugin.so

events:
  - recipient: true
    regex: ".*hello"
    #response: Greetings!
    plugin: Example
    type: message_event
    options:
      lookup_user: true
```