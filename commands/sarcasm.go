package commands

import (
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
)

// Sarcasm sarcasm looool
func Sarcasm(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "</sarcasm>")
	bot.Send(msg)

	_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.MessageID,
	})

	if err != nil {
		common.Log.Warn(err.Error())
	}
}
