package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "</sarcasm>")
	api.Send(msg)

	_, _ = api.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
}

func init() {
	commands.Register("s", &Command{})
}
