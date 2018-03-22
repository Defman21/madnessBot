package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

// Reverse lul
func Reverse(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := update.Message.ReplyToMessage.Text
	r := []rune(text)

	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(r))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
