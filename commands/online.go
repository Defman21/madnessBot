package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"madnessBot/state/online"
	"madnessBot/templates"
)

type OnlineCmd struct{}

func (c OnlineCmd) UseLua() bool {
	return false
}

func (c OnlineCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	streamers := online.GetOnline()
	if len(streamers) == 0 {
		helpers.SendMessage(api, update, "Никто не стримит", false, true)
		return
	}
	msg := templates.ExecuteTemplate("commands_online", streamers)
	helpers.SendMessage(api, update, msg, false, false)
}
