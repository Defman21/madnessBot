package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"time"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	photo.FileID = fmt.Sprintf("https://thecatapi.com/api/images/get?type=jpg,png&%s", timestamp)
	photo.UseExisting = true

	_, err := api.Send(photo)
	if err != nil {
		msg := fmt.Sprintf("Кот не нашелбся....")
		_, _ = api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	}

}

func init() {
	commands.Register("cat", &Command{})
}
