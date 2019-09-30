package commands

import (
	"github.com/Defman21/madnessBot/common"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
	"strings"
)

func PayCheck(bot *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	payingUsers := map[int]bool{}
	for _, user := range strings.Split(os.Getenv("PAYING_USERS"), ";") {
		userID, err := strconv.Atoi(user)
		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to convert userID str to int")
			continue
		}
		payingUsers[userID] = true
	}

	user := update.Message.From.ID
	if _, exists := payingUsers[user]; exists {
		return true
	}

	msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID,
		"AwADAgADwgADC6ZpS13yfdzm_pTzAg")
	bot.Send(msg)
	return false
}
