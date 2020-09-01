package main

import (
	"fmt"
	"time"
	"github.com/BurntSushi/toml"
	"github.com/robfig/cron/v3"
  "github.com/nlopes/slack"
)

type tomlConfig struct {
	SlackToken string `toml:"slack_token"`
	TargetChannels map[string]targetChannel `toml:"target_channels"`
}

type targetChannel struct {
	ChannelId string `toml:"channel_id"`
	Timer string `toml:"timer"`
}



func cleanChats(slackApi *slack.Client, config *tomlConfig) {
	apiCallCounter := 0
}

func main () {

	// load config
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println("configuration file load failed.")
		return
	}

	slackApi := slack.New(config.SlackToken)
	cronInstance := cron.New()

	cronInstance.AddFunc("@every 1m", cleanChats(slackApi, config))

	for {
		time.Sleep(86400 * time.Second)
	}

}