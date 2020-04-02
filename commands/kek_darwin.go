package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
)

type KekCmd struct{}

func (c KekCmd) UseLua() bool {
	return false
}

func (c KekCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !helpers.PayCheck(api, update) {
		return
	}
}
