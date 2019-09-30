package common

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
	"strings"
)

// IsAdmin 4Head
func IsAdmin(user *tgbotapi.User) bool {
	admins := map[int]bool{}
	for _, admin := range strings.Split(os.Getenv("ADMINS_LIST"), ";") {
		adminID, err := strconv.Atoi(admin)
		if err != nil {
			Log.Error().Err(err).Msg("Failed to convert admin ID from str to int")
		}
		admins[adminID] = true
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
		Log.Error().Err(err).Msg("Failed to get chat administrators")
		return false
	}

	for _, chatMember := range admins {
		if chatMember.User.ID == update.Message.From.ID {
			return true
		}
	}
	return false
}
