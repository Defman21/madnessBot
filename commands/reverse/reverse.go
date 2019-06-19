package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"gopkg.in/telegram-bot-api.v4"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return true
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := update.Message.ReplyToMessage.Text
	r := []rune(text)

	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(r))
	msg.ReplyToMessageID = update.Message.MessageID

	api.Send(msg)
}

func init() {
	commands.Register("reverse", &Command{})
}
