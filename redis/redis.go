package redis

import (
	"github.com/go-redis/redis/v8"
	"madnessBot/common/logger"
	"madnessBot/config"
)

var instance *redis.Client

func Init() {
	instance = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	logger.Log.Info().Msg("Initialized Redis")
}

func Get() *redis.Client {
	return instance
}
