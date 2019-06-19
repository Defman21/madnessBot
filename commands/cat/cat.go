package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	placeholder := tgbotapi.NewMessage(update.Message.Chat.ID, "ищу котека...")
	placeholderMessage, _ := api.Send(placeholder)

	photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	photo.FileID = fmt.Sprintf("https://thecatapi.com/api/images/get?type=jpg,png&%s", timestamp)
	photo.UseExisting = true

	_, err := api.Send(photo)
	if err != nil {
		msg := fmt.Sprintf("Кот не нашелбся....")
		_, _ = api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	}

	_, _ = api.DeleteMessage(tgbotapi.DeleteMessageConfig{
		MessageID: placeholderMessage.MessageID,
		ChatID:    placeholderMessage.Chat.ID,
	})

}

func init() {
	commands.Register("cat", &Command{})
}
