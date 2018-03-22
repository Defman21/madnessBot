package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

// Kek zulul
func Kek(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !payCheck(bot, update) {
		return
	}
}
