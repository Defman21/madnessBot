package boosty

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"strconv"
)

const redisKey = "madnessBot:state:payers"

var log = &logger.Log

func HandleUpdate(_ *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var _redis = redis.Get()

	if len(update.Message.NewChatMembers) > 0 {
		for _, user := range update.Message.NewChatMembers {
			userId := strconv.FormatInt(int64(user.ID), 10)
			_, err := _redis.HSet(redisKey, userId, true).Result()
			if err != nil {
				log.Error().Err(err).
					Str("key", redisKey).
					Str("value", userId).
					Msg("Failed to HSET redis key")
				continue
			}
			log.Info().Int("user-id", user.ID).Msg("Added to payers")
		}
	}

	if update.Message.LeftChatMember != nil {
		userId := strconv.FormatInt(int64(update.Message.LeftChatMember.ID), 10)
		_, err := _redis.HDel(redisKey, userId).Result()
		if err != nil {
			log.Error().Err(err).
				Str("key", redisKey).
				Str("value", userId).
				Msg("Failed to HDEL redis key")
			return
		}
		log.Info().Int("user-id", update.Message.LeftChatMember.ID).Msg("Removed from payers")
	}
}

func GetPayers() map[int64]bool {
	res := map[int64]bool{}
	usersMap, err := redis.Get().HGetAll(redisKey).Result()
	if err != nil {
		log.Error().Err(err).Str("key", redisKey).Msg("Failed to HGETALL redis key")
	}

	for k := range usersMap {
		i, _ := strconv.ParseInt(k, 10, 64)
		res[i] = true
	}
	return res
}
