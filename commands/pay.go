package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

var payers, notified map[int]bool

func init() {
	payers = map[int]bool{323141608: true, 306022838: true, 71524437: true}
	notified = make(map[int]bool)
}

func payCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	user := update.Message.From.ID
	if _, exists := payers[user]; exists {
		return true
	}
	if _, exists := notified[user]; !exists {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "PLATI NOLOGI")
		bot.Send(msg)
		notified[user] = true
	}
	return false
}
