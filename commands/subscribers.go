package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"madnessBot/templates"
)

type SubscribersCmd struct{}
type subscribeUsers map[string]string

func (c SubscribersCmd) UseLua() bool {
	return false
}

func getSubscribersList() (users subscribeUsers) {
	users, err := redis.Get().HGetAll(redisKey).Result()

	if err != nil {
		logger.Log.Error().Err(err).Str("key", redisKey).Msg("Failed to HGETALL redis key")
		return
	}

	return
}

func (c SubscribersCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	users := getSubscribersList()
	if users == nil {
		logger.Log.Warn().Msg("Empty user list")
		return
	}

	subscribers := templates.ExecuteTemplate("commands_subscribers", users)

	helpers.SendMessage(api, update, subscribers, true, true)
}