package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"math/rand"
	"strings"
	"time"
)

type ResolveCmd struct{}

func (c ResolveCmd) UseLua() bool {
	return true
}

func (c ResolveCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := update.Message.CommandArguments()
	choices := strings.Split(text, "/")
	var choice string

	if len(text) == 0 {
		choice = "и че я должен выбрать? forsenThink"
	} else if len(choices) == 1 {
		choice = "хочу больше вариантов"
	} else {
		source := rand.NewSource(time.Now().Unix())
		random := rand.New(source)
		choice = choices[random.Intn(len(choices))]
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpers.EscapeMarkdownV2(choice))
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.ReplyToMessageID = update.Message.MessageID

	_, _ = api.Send(msg)
}
