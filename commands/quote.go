package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"math/rand"
	"strconv"
)

type quoteFile struct {
	Quotes []struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	} `json:"quotes"`
}

var Quotes quoteFile

func init() {
	str, err := ioutil.ReadFile("./data/quotes.json")
	if err != nil {
		common.Log.Warn("Failed to read ./data/quotes.json")
		return
	}

	json.Unmarshal(str, &Quotes)

	common.Log.WithFields(logrus.Fields{
		"bytes": string(str),
		"quote": Quotes,
	}).Debug("Dump")
}

func Quote(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	quoteID, err := strconv.Atoi(update.Message.CommandArguments())

	if err != nil {
		common.Log.Warn("Not a valid number")
		quoteID = rand.Intn(len(Quotes.Quotes) - 1)
	}

	for _, quote := range Quotes.Quotes {
		if quote.ID == quoteID {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Цитатка #%d: %s", quote.ID, quote.Text))
			bot.Send(msg)
			return
		}
	}
	common.Log.Warn("Quote not found")
}
