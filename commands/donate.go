package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"madnessBot/config"
)

type DonateCmd struct{}

func (c DonateCmd) UseLua() bool {
	return true
}

func (c DonateCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	helpers.SendMessage(api, update, config.Config.BoostyLink, true, true)
}
