package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"gopkg.in/telegram-bot-api.v4"
)

// Cat zulul
func Cat(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !payCheck(bot, update) {
		return
	}
	res, err := goreq.Request{
		Uri: "https://thecatapi.com/api/images/get",
		QueryString: struct {
			Type string
		}{
			Type: "jpg,png",
		},
	}.Do()

	if err != nil {
		common.Log.Warn(err)
		return
	}
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
	msg.FileID = res.Header.Get("Location")
	msg.UseExisting = true
	_, err = bot.Send(msg)
	if err != nil {
		msg := fmt.Sprintf("Все летит в пизду\n%s\nURL: %s", err, res.Header.Get("Location"))
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	}
}
