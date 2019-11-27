package common

import (
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// IsAdmin 4Head
func IsAdmin(user *tgbotapi.User) bool {
	admins := map[int]bool{}
	for _, admin := range config.Config.Admins {
		admins[admin] = true
	}

	if _, exists := admins[user.ID]; exists {
		return true
	}
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
