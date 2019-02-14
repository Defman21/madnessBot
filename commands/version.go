package commands

import (
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
	"os/exec"
)

// Version reports last commit
func Version(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		return
	}

	commit, _ := exec.Command("git", "show", "-s").Output()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(commit))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
