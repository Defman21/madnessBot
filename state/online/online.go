package online

import (
	"context"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"strconv"
)

var log = &logger.Log

func Add(username string, isOnline bool) {
	redis.Get().HSet(context.Background(), redis.OnlineStreamersKey, username, isOnline)
}

func GetOnline() (result []string) {
	kvPair, err := redis.Get().HGetAll(context.Background(), redis.OnlineStreamersKey).Result()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get redis key %s", redis.OnlineStreamersKey)
	}
	for username, isOnlineStr := range kvPair {
		isOnline, _ := strconv.ParseBool(isOnlineStr)
		if isOnline {
			result = append(result, username)
		}
	}
	return
}
