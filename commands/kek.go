package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
)

func kek(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	photos, err := bot.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(update.Message.From.ID))
	if err != nil {
		common.Log.Warn("LUL ZULUL")
	} else {
		fmt.Printf("%+v", photos)
	}
}
