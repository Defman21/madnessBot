package commands

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/redis"
)

const redisKey = "madnessBot:state:subscriptions"

type SubscribeCmd struct{}

func (c SubscribeCmd) UseLua() bool {
	return false
}

func (c SubscribeCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	channel := update.Message.CommandArguments()
	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}
	userID, found := helpers.GetTwitchUserIDByLogin(channel)
	if found {
		if err := helpers.SendEventSubMessage(channel, helix.EventSubTypeStreamOnline); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to subscribe")
			return
		}

		_, err := redis.Get().HSet(context.Background(), redisKey, channel, userID).Result()

		if err != nil {
			log.Error().Err(err).
				Str("key", redisKey).
				Str("field", channel).
				Str("value", userID).
				Msg("Failed to HSET redis key")
			return
		}

		helpers.SendMessage(
			api,
			update,
			helpers.EscapeMarkdownV2(fmt.Sprintf("Бот теперь аки маньяк будет преследовать %s "+
				"до конца своих дней.", channel)), true, true,
		)
	}
}
