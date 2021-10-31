package resubscribe

import (
	"context"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"time"
)

const expireTime = time.Hour * 24 * 7

func SaveState() {
	_, err := redis.Get().Set(context.Background(), redis.ResubscribeAtKey, true, expireTime).Result()

	if err != nil {
		logger.Log.Error().Err(err).
			Str("key", redis.ResubscribeAtKey).
			Bool("value", true).
			Dur("ex", expireTime).
			Msg("Failed to SET redis key")
		return
	}
}

func GetState() *time.Time {
	timestamp, err := redis.Get().TTL(context.Background(), redis.ResubscribeAtKey).Result()

	if err != nil {
		logger.Log.Error().Err(err).Str("key", redis.ResubscribeAtKey).Msg("Failed to TTL redis key")
		return nil
	}

	if timestamp <= 0 {
		logger.Log.Warn().Str("key", redis.ResubscribeAtKey).Msg("Redis key expired")
		timestamp = -1 * time.Second
	}

	after := time.Now().Local().Add(timestamp)
	return &after
}
