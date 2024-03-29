package commands

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicklaw5/helix/v2"
	"madnessBot/common"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/redis"
)

type UnsubscribeCmd struct{}

func (c UnsubscribeCmd) UseLua() bool {
	return false
}

func generateUnsubscribeTopic(userID string) string {
	return fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID)
}

func (c UnsubscribeCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		helpers.SendMessage(api, update, "TriHard LULW", true, true)
		return
	}

	channel := update.Message.CommandArguments()

	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}

	users := getSubscribersList()

	if userID, ok := users[channel]; ok {
		go func(channel string, userID string) {
			if err := helpers.UnsubscribeFromEventSub(channel, helix.EventSubTypeStreamOnline); err != nil {
				logger.Log.Error().Err(err).Msg("Failed to send a request")
				return
			}

			logger.Log.Info().Str("user", channel).Msg("Unsubscribed")

			_, err := redis.Get().HDel(context.Background(), redis.SubscriptionsKey, channel).Result()
			if err != nil {
				logger.Log.Error().Err(err).
					Str("key", redis.SubscriptionsKey).
					Str("field", channel).
					Msg("Failed to HDEL redis key")
			}

			helpers.SendMessage(api, update, fmt.Sprintf("Unsubscribed from %s", channel), true, true)
		}(channel, userID)
	} else {
		logger.Log.Warn().Str("channel", channel).Msg("Channel not found")
	}
}
