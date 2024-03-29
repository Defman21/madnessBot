package notify

import (
	"context"
	"fmt"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"strings"
)

var _redis = redis.Get

var log = &logger.Log

func Add(userID string, userName string) {
	_redis().RPush(context.Background(), fmt.Sprintf(redis.NotificationsKey, userID), userName)
}

func Remove(userID string, userName string) {
	_redis().LRem(context.Background(), fmt.Sprintf(redis.NotificationsKey, userID), 1, userName)
}

func GenerateString(userID string) string {
	redisKey := getRedisKey(userID)
	userLogins, err := _redis().LRange(context.Background(), redisKey, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get range from redis key %s", redisKey)
		return ""
	}
	return strings.Join(userLogins, ", ")
}

func getRedisKey(userID string) string {
	return fmt.Sprintf(redis.NotificationsKey, userID)
}
