package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/notifier"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	login := update.Message.CommandArguments()
	userName := update.Message.From.UserName
	userID, found := common.GetTwitchUserIDByLogin(login)
	if !found {
		common.SendMessage(api, update.Message.Chat.ID, "стример не найден", nil)
		return
	}

	notifier.Get().Add(userID, fmt.Sprintf("@%s", userName))

	common.SendMessage(
		api,
		update.Message.Chat.ID,
		fmt.Sprintf("подписал тебя на оповещения от %s", login),
		&update.Message.MessageID,
	)
}

func init() {
	commands.Register("notify_me", &Command{})
}
