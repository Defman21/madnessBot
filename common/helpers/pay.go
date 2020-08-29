package helpers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common"
	"madnessBot/integrations/boosty"
)

func PayCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	user := int64(update.Message.From.ID)
	if common.IsAdmin(update.Message.From) {
		return true
	}
	if _, exists := boosty.GetPayers()[user]; exists {
		return true
	}

	SendInvalidArgumentsMessage(bot, update)
	return false
}
