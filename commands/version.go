package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common"
	"madnessBot/common/helpers"
	"os/exec"
)

type VersionCmd struct{}

func (c VersionCmd) UseLua() bool {
	return true
}

func (c VersionCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		return
	}

	commit, _ := exec.Command("git", "show", "-s").Output()
	helpers.SendMessage(api, update, string(commit), true, false)
}
