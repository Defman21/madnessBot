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
		helpers.SendMessage(api, update, "TriHard LULW", true)
		return
	}

	channel := update.Message.CommandArguments()

	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}

	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to read users.json")
		return
	}

	var users Users

	_ = json.Unmarshal(bytes, &users)

	if userID, ok := users[channel]; ok {
		go func(channel string, userID string) {
			if errs := helpers.SendTwitchHubMessage(channel, "unsubscribe", generateTopic(userID)); errs != nil {
				logger.Log.Error().Errs("errs", errs).Msg("Failed to send a request")
				return
			}

			logger.Log.Info().Str("user", channel).Msg("Unsubscribed")

			delete(users, channel)

			jsonStr, _ := json.Marshal(users)
			err = ioutil.WriteFile("./data/users.json", []byte(jsonStr), 0644)

			if err == nil {
				logger.Log.Info().Msg("Updated users.json")
				helpers.SendMessage(
					api,
					update,
					fmt.Sprintf("Unsubscribed from %s", channel),
					true,
				)
			} else {
				logger.Log.Warn().Err(err).Msg("Couldn't write to users.json")
			}
		}(channel, userID)
	} else {
		logger.Log.Warn().Str("channel", channel).Msg("Channel not found")
	}
}

func init() {
	commands.Register("unsubscribe", &Command{})
}
