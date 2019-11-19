package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
)

type Command struct{}
type Users map[string]string

func (c *Command) UseLua() bool {
	return false
}

func generateTopic(userID string) string {
	return fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID)
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		return
	}
	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to read users.json")
		return
	}

	var users Users

	_ = json.Unmarshal(bytes, &users)

	for channel, userID := range users {
		go func(channel string, userID string) {
			if errs := helpers.SendTwitchHubMessage(channel, "subscribe", generateTopic(userID)); errs != nil {
				logger.Log.Error().Errs("errs", errs).Msg("Failed to resubscribe")
			} else {
				logger.Log.Info().Str("channel", channel).Msg("Subscribed")
			}
		}(channel, userID)
	}

	common.ResubscribeState.Save()
}

func init() {
	commands.Register("resubscribe", &Command{})
}
