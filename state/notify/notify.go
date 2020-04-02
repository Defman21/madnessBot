package notify

import (
	"fmt"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"strings"
)

const redisKey = "madnessBot:state:notify:%s"

var _redis = redis.Get

var log = &logger.Log

func Add(userID string, userName string) {
	_redis().RPush(fmt.Sprintf(redisKey, userID), userName)
}

func Remove(userID string, userName string) {
	_redis().LRem(fmt.Sprintf(redisKey, userID), 1, userName)
}

func GenerateString(userID string) string {
	redisKey := getRedisKey(userID)
	userLogins, err := _redis().LRange(redisKey, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get range from redis key %s", redisKey)
		return ""
	}
	return strings.Join(userLogins, ", ")
}

func getRedisKey(userID string) string {
	return fmt.Sprintf(redisKey, userID)
}
