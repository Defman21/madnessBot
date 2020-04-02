package online

import (
	"madnessBot/common/logger"
	"madnessBot/redis"
	"strconv"
)

const redisKey = "madnessBot:state:online"

var log = &logger.Log

func Add(username string, isOnline bool) {
	redis.Get().HSet(redisKey, username, isOnline)
}

func GetOnline() (result []string) {
	kvPair, err := redis.Get().HGetAll(redisKey).Result()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get redis key %s", redisKey)
	}
	for username, isOnlineStr := range kvPair {
		isOnline, _ := strconv.ParseBool(isOnlineStr)
		if isOnline {
			result = append(result, username)
		}
	}
	return
}
