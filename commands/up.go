package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

func Up(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewStickerShare(update.Message.Chat.ID,
		"CAADAgADPgAD43TSFv8rTPYvm_MJAg")
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
