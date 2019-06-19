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
	msg := tgbotapi.NewStickerShare(update.Message.Chat.ID,
		"CAADAgADPgAD43TSFv8rTPYvm_MJAg")
	msg.ReplyToMessageID = update.Message.MessageID

	api.Send(msg)
}

func init() {
	commands.Register("up", &Command{})
}
