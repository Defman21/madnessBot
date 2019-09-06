package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strings"
	"time"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return true
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := update.Message.CommandArguments()
	choices := strings.Split(text, "/")
	var choice string

	if len(text) == 0 {
		choice = "и че я должен выбрать? forsenThink"
	} else if len(choices) == 1 {
		choice = "хочу больше вариантов"
	} else {
		source := rand.NewSource(time.Now().Unix())
		random := rand.New(source)
		choice = choices[random.Intn(len(choices))]
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, choice)
	msg.ReplyToMessageID = update.Message.MessageID

	_, _ = api.Send(msg)
}

func init() {
	commands.Register("resolve", &Command{})
	commands.Register("r", &Command{})
}
