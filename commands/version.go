package commands

import (
	"gopkg.in/telegram-bot-api.v4"
	"os/exec"
)

// Version reports last commit
func Version(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	commit, _ := exec.Command("git", "show", "-s").Output()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(commit))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
