package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
)

func AddQuote(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	text := update.Message.CommandArguments()

	ID := Quotes.Quotes[len(Quotes.Quotes)-1].ID + 1

	Quotes.Quotes = append(Quotes.Quotes, struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	}{
		ID:   ID,
		Text: text,
	})

	bytes, err := json.Marshal(Quotes)
	common.Log.WithFields(logrus.Fields{
		"quote": Quotes,
		"json":  string(bytes),
	}).Debug("Dump")

	if err != nil {
		common.Log.Warn("Failed to serialize the object")
	} else {
		err := ioutil.WriteFile("./data/quotes.json", bytes, 0644)
		if err != nil {
			common.Log.Warn("Failed to write to the file")
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				fmt.Sprintf("Добавлена цитатка #%d: %s", ID, text))
			bot.Send(msg)
		}
	}
}
