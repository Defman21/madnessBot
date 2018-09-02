package commands

import (
	"fmt"
	"strconv"
	"time"
	"gopkg.in/telegram-bot-api.v4"
)

// Cat zulul
func Cat(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msg.FileID = fmt.Sprintf("https://thecatapi.com/api/images/get?type=jpg,png&%s", timestamp)
	msg.UseExisting = true
	_, err := bot.Send(msg)
	if err != nil {
		msg := fmt.Sprintf("Все летит в пизду\n%s\nURL: %s", err, msg.FileID)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	}
}
