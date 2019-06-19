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
	placeholder := tgbotapi.NewPhotoShare(
		update.Message.Chat.ID,
		"https://static.thenounproject.com/png/101791-200.png",
	)
	placeholder.Caption = "ищу котека..."
	placeholderMessage, _ := api.Send(placeholder)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	editmsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    placeholderMessage.Chat.ID,
			MessageID: placeholderMessage.MessageID,
		},
		Media: tgbotapi.BaseInputMedia{
			Type:    "photo",
			Media:   fmt.Sprintf("https://thecatapi.com/api/images/get?type=jpg,png&%s", timestamp),
			Caption: "котек найден!",
		},
	}

	_, err := api.Send(editmsg)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to edit a message")
	}
}

func init() {
	commands.Register("cat", &Command{})
}
