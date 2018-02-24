package commands

import (
	"gopkg.in/telegram-bot-api.v4"
	"os"
)

func Birthday(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, os.Getenv("BIRTHDAY_URL"))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
