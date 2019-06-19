package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var name string
	if update.Message.From.LastName != "" {
		name = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
	} else {
		name = update.Message.From.FirstName
	}
	text := update.Message.CommandArguments()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("_%s %s_", name, text))
	msg.ParseMode = tgbotapi.ModeMarkdown

	api.Send(msg)

	_, err := api.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.MessageID,
	})

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to send a message")
	}
}

func init() {
	commands.Register("me", &Command{})
}
