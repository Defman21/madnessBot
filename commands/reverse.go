package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
)

type ReverseCmd struct{}

func (c ReverseCmd) UseLua() bool {
	return true
}

func (c ReverseCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.Message.ReplyToMessage == nil {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}
	text := update.Message.ReplyToMessage.Text
	r := []rune(text)

	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	helpers.SendMessage(api, update, helpers.EscapeMarkdownV2(string(r)), true, false)
}
