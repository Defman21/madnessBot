package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"madnessBot/state/notify"
)

type UnnotifyMeCmd struct{}

func (c UnnotifyMeCmd) UseLua() bool {
	return false
}

func (c UnnotifyMeCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	login := update.Message.CommandArguments()
	userName := update.Message.From.UserName
	userID, found := helpers.GetTwitchUserIDByLogin(login)

	if !found {
		helpers.SendMessage(api, update, "стример не найден", true, true)
		return
	}

	notify.Remove(userID, fmt.Sprintf("@%s", userName))

	helpers.SendMessage(api, update, fmt.Sprintf("отписал тебя от оповещений %s", login), true, true)
}
