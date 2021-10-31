package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicklaw5/helix/v2"
	"madnessBot/common"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/state/resubscribe"
)

type ResubscribeCmd struct{}

func (c ResubscribeCmd) UseLua() bool {
	return false
}

func generateTopic(userID string) string {
	return fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID)
}

func (c ResubscribeCmd) Run(_ *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		return
	}

	users := getSubscribersList()

	for channel, userID := range users {
		go func(channel string, userID string) {
			if err := helpers.SendEventSubMessage(channel, helix.EventSubTypeStreamOnline); err != nil {
				logger.Log.Error().Err(err).Msg("Failed to resubscribe")
			} else {
				logger.Log.Info().Str("channel", channel).Msg("Subscribed")
			}
		}(channel, userID)
	}

	resubscribe.SaveState()
}
