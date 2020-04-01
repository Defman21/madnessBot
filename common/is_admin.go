package common

import (
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// IsAdmin 4Head
func IsAdmin(user *tgbotapi.User) bool {
	if _, exists := config.Config.GetAdmins()[int64(user.ID)]; exists {
		return true
	}
	logger.Log.Warn().
		Interface("admins", config.Config.Admins).
		Int("user_id", user.ID).
		Msg("Not an admin")
	return false
}

// IsMod 4Head
func IsMod(api *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	admins, err := api.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: update.Message.Chat.ChatConfig(),
	})
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get chat administrators")
		return false
	}

	for _, chatMember := range admins {
		if chatMember.User.ID == update.Message.From.ID {
			return true
		}
	}
	return false
}
