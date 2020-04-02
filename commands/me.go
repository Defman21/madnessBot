package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
)

type MeCmd struct{}

func (c MeCmd) UseLua() bool {
	return false
}

func (c MeCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var name string
	if update.Message.From.LastName != "" {
		name = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
	} else {
		name = update.Message.From.FirstName
	}
	text := update.Message.CommandArguments()
	helpers.SendMessage(api, update, fmt.Sprintf("_%s %s_", name, text), false, false)
	_, _ = api.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
}
