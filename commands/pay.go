package commands

import (
	"github.com/Defman21/madnessBot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PayCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	user := int64(update.Message.From.ID)
	if _, exists := config.Config.GetPayers()[user]; exists {
		return true
	}

	msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID, "AwADAgADwgADC6ZpS13yfdzm_pTzAg")
	bot.Send(msg)
	return false
}
