package commands

import (
	"github.com/Defman21/madnessBot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PayCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	payingUsers := map[int]bool{}
	for user := range config.Config.Payers {
		payingUsers[user] = true
	}

	user := update.Message.From.ID
	if _, exists := payingUsers[user]; exists {
		return true
	}

	msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID, "AwADAgADwgADC6ZpS13yfdzm_pTzAg")
	bot.Send(msg)
	return false
}
