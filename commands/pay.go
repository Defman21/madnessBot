package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var payers, notified map[int]bool

func init() {
	payers = map[int]bool{
		323141608: true, // Tishka
		306022838: true, // Kleozis
		71524437:  true, // defman
		105513756: true, // defman
		//301864265: true, // borobushe
		//431674591: true, // Refferency
		//86097149: true // advancher
	}
	notified = make(map[int]bool)
}

func PayCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	user := update.Message.From.ID
	if _, exists := payers[user]; exists {
		return true
	}
	if _, exists := notified[user]; !exists {
		msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID,
			"AwADAgADwgADC6ZpS13yfdzm_pTzAg")
		bot.Send(msg)
		notified[user] = true
	}
	return false
}
