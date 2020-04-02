package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/state/notify"
)

type NotifyMeCmd struct{}

func (c NotifyMeCmd) UseLua() bool {
	return false
}

func (c NotifyMeCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	login := update.Message.CommandArguments()
	existingSubscribers := getSubscribersList()

	if existingSubscribers == nil {
		logger.Log.Warn().Msg("Empty user list")
		return
	}

	if _, exists := existingSubscribers[login]; !exists {
		helpers.SendMessage(api, update, "бот не подписан на этого юзера", true, true)
		return
	}

	userName := update.Message.From.UserName
	userID, found := helpers.GetTwitchUserIDByLogin(login)
	if !found {
		helpers.SendMessage(api, update, "стример не найден", true, true)
		return
	}

	notify.Add(userID, fmt.Sprintf("@%s", userName))

	helpers.SendMessage(api, update, fmt.Sprintf("подписал тебя на оповещения от %s", login), true, true)
}
