package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "</sarcasm>")
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
	commands.Register("sacrasm", &Command{})
}
