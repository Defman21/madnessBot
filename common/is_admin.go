package common

import (
	"gopkg.in/telegram-bot-api.v4"
)

// IsAdmin 4Head
func IsAdmin(user *tgbotapi.User) bool {
	var admins = map[int]bool{
		71524437: true,
		105513756: true,
	}
	if _, exists := admins[user.ID]; exists {
		return true
	}
	return false
}
