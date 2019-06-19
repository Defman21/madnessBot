package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os/exec"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return true
}

func (c *Command) Run(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		return
	}

	commit, _ := exec.Command("git", "show", "-s").Output()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(commit))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}

func init() {
	commands.Register("version", &Command{})
}
