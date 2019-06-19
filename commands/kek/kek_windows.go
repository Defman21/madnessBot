package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"gopkg.in/telegram-bot-api.v4"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) RunGo(api *tgbotapi.BotAPI, update *tgbotapi.Update) {}

func init() {
	commands.Register("kek", &Command{})
}
