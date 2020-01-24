package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/commands"
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
	channel := update.Message.CommandArguments()
	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}
	userID, found := helpers.GetTwitchUserIDByLogin(channel)
	if found {
		if errs := helpers.SendTwitchHubMessage(channel, "subscribe", generateTopic(userID)); errs != nil {
			logger.Log.Error().Errs("errs", errs).Msg("Failed to subscribe")
			return
		}

		var users Users
		bytes, err := ioutil.ReadFile("./data/users.json")
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to read users.json")
		}

		json.Unmarshal(bytes, &users)

		users[channel] = userID
		bytes, err = json.Marshal(users)

		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to serialize users")
		} else {
			err := ioutil.WriteFile("./data/users.json", bytes, 0644)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Failed to write users.json")
				return
			}
			helpers.SendMessage(api, update, fmt.Sprintf("Бот теперь аки маньяк будет преследовать %s "+
				"до конца своих дней.", channel), true, true)
		}
	}
}

func init() {
	commands.Register("subscribe", &Command{})
}
