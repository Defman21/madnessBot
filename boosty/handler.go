package boosty

import (
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var log = &logger.Log

func HandleUpdate(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if len(update.Message.NewChatMembers) > 0 {
		for _, user := range update.Message.NewChatMembers {
			config.Config.AddPayer(int64(user.ID))
			log.Info().Int("user-id", user.ID).Msg("Added to payers")
		}
	}

	if update.Message.LeftChatMember != nil {
		config.Config.RemovePayer(int64(update.Message.LeftChatMember.ID))
		log.Info().Int("user-id", update.Message.LeftChatMember.ID).Msg("Removed from payers")
	}

	config.Save()
}
