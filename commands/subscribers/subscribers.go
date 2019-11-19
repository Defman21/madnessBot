package commands

import (
	"encoding/json"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/templates"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
)

type Command struct{}
type Users map[string]string

func (c *Command) UseLua() bool {
	return false
}

func GetList() (users Users) {
	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to read users.json")
		return nil
	}

	json.Unmarshal(bytes, &users)

	return users
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	users := GetList()
	if users == nil {
		logger.Log.Warn().Msg("Empty user list")
		return
	}

	subscribers := templates.ExecuteTemplate("commands_subscribers", users)

	helpers.SendMessage(api, update, subscribers, true)
}

func init() {
	commands.Register("subs", &Command{})
}
