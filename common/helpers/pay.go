package helpers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/integrations/boosty"
)

func PayCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	user := int64(update.Message.From.ID)
	if _, exists := boosty.GetPayers()[user]; exists {
		return true
	}

	SendInvalidArgumentsMessage(bot, update)
	return false
}
