package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/online_state"
	"github.com/Defman21/madnessBot/templates"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	online := online_state.GetOnline()
	if len(online) == 0 {
		helpers.SendMessage(api, update, "Никто не стримит", false)
		return
	}
	msg := templates.ExecuteTemplate("commands_online", online)
	helpers.SendMessage(api, update, msg, false)
}

func init() {
	commands.Register("online", &Command{})
}
