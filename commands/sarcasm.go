package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
)

type SarcasmCmd struct{}

func (c SarcasmCmd) UseLua() bool {
	return false
}

func (c SarcasmCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	helpers.SendMessage(api, update, "</sarcasm>", false, false)
	_, _ = api.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
}
