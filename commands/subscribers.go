package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
)

// Subscribers list subscribers
func Subscribers(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to read users.json")
		return
	}

	var users Users
	var subscribers string

	json.Unmarshal(bytes, &users)

	for username, _ := range users {
		subscribers += fmt.Sprintf("%s\n", username)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, subscribers)
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
