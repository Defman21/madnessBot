package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
)

type Command struct{}
type Users map[string]string

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
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

	api.Send(msg)
}

func init() {
	commands.Register("subs", &Command{})
}
