package commands

import (
	"gopkg.in/telegram-bot-api.v4"
	"os"
)

func Donate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, os.Getenv("CARD_NUMBER"))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
