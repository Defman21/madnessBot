package resubscribe

import (
	"context"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"time"
)

const redisKey = "madnessBot:state:subscriptions:resubscribeAt"
const expireTime = time.Hour * 24 * 7

func SaveState() {
	_, err := redis.Get().Set(context.Background(), redisKey, true, expireTime).Result()

	if err != nil {
		logger.Log.Error().Err(err).
			Str("key", redisKey).
			Bool("value", true).
			Dur("ex", expireTime).
			Msg("Failed to SET redis key")
		return
	}
}

func GetState() *time.Time {
	timestamp, err := redis.Get().TTL(context.Background(), redisKey).Result()

	if err != nil {
		logger.Log.Error().Err(err).Str("key", redisKey).Msg("Failed to TTL redis key")
		return nil
	}

	if timestamp <= 0 {
		logger.Log.Warn().Str("key", redisKey).Msg("Redis key expired")
		timestamp = -1 * time.Second
	}

	after := time.Now().Local().Add(timestamp)
	return &after
}
