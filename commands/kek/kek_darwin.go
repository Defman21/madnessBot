package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {}

func init() {
	commands.Register("kek", &Command{})
}
