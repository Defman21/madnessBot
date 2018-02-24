package commands

import (
	"bytes"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
)

func QuoteList(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var buffer bytes.Buffer
	for _, quote := range Quotes.Quotes {
		str := fmt.Sprintf("%d. %s\n", quote.ID, quote.Text)
		buffer.WriteString(str)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, buffer.String())
	bot.Send(msg)
}
