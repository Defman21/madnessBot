package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"gopkg.in/telegram-bot-api.v4"
	"os"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return true
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, os.Getenv("CARD_NUMBER"))
	msg.ReplyToMessageID = update.Message.MessageID

	api.Send(msg)
}

func init() {
	commands.Register("donate", &Command{})
}
