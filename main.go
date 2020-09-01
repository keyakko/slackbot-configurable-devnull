package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/robfig/cron/v3"
	"github.com/slack-go/slack"
)

type tomlConfig struct {
	SlackToken     string                   `toml:"slack_token"`
	TargetChannels map[string]targetChannel `toml:"target_channels"`
}

type targetChannel struct {
	ChannelID string `toml:"channel_id"`
	Timer     int64  `toml:"timer"`
}

func cleanChats(slackClient *slack.Client, config tomlConfig) {
	// apiCallCounter := 0
	currentTime := time.Now().Unix()

	for _, targetChannelConf := range config.TargetChannels {
		fmt.Println(targetChannelConf.ChannelID)
		fmt.Println(targetChannelConf.Timer)

		var getConversasionHistoryParameters slack.GetConversationHistoryParameters
		getConversasionHistoryParameters.ChannelID = targetChannelConf.ChannelID
		getConversasionHistoryParameters.Latest = strconv.FormatInt(currentTime-targetChannelConf.Timer, 10)
		getConversasionHistoryParameters.Cursor = ""
		getConversasionHistoryParameters.Inclusive = false
		getConversasionHistoryParameters.Limit = 100
		getConversasionHistoryParameters.Oldest = "0"

		conversationHistory, err := slackClient.GetConversationHistory(
			&getConversasionHistoryParameters,
		)
		// apiCallCounter++
		if err != nil {
			fmt.Println("[ERROR] getConversationHistory failed.")
			fmt.Println(err)
			continue
		}

		for _, message := range conversationHistory.Messages {
			slackClient.DeleteMessage(targetChannelConf.ChannelID, message.Timestamp)
		}

	}
}

func main() {

	// load config
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println("configuration file load failed.")
		return
	}

	slackClient := slack.New(config.SlackToken)
	cronInstance := cron.New()

	cronInstance.AddFunc("@every 1m", func() { cleanChats(slackClient, config) })
	cronInstance.Start()
	cleanChats(slackClient, config)

	for {
		time.Sleep(86400 * time.Second)
	}

}
