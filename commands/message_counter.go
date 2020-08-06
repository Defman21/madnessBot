package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/redis"
	"madnessBot/templates"
)

type MessageCounterCmd struct{}

func (c MessageCounterCmd) UseLua() bool {
	return true
}

type messageCounterCommandTemplate struct {
	Count string
}

func (c MessageCounterCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	counterStr, err := redis.Get().Get("madnessBot:messageCounter").Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get message counter")
		return
	}
	msg := templates.ExecuteTemplate("commands_message_counter",
		messageCounterCommandTemplate{
			Count: counterStr,
		},
	)
	helpers.SendMessage(api, update, msg, true, true)
}
