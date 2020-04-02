package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func generateSubscribeTopic(userID string) string {
	return fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID)
}

func (c SubscribeCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	channel := update.Message.CommandArguments()
	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}
	userID, found := helpers.GetTwitchUserIDByLogin(channel)
	if found {
		if errs := helpers.SendTwitchHubMessage(channel, "subscribe", generateSubscribeTopic(userID)); errs != nil {
			logger.Log.Error().Errs("errs", errs).Msg("Failed to subscribe")
			return
		}

		_, err := redis.Get().HSet(redisKey, channel, userID).Result()

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
			fmt.Sprintf("Бот теперь аки маньяк будет преследовать %s "+
				"до конца своих дней.", channel), true, true,
		)
	}
}
