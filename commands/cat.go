package commands

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marpaia/graphite-golang"
	"madnessBot/common/logger"
	"madnessBot/common/metrics"
	"math/rand"
	"strconv"
	"time"
)

type CatCmd struct{}

func (c CatCmd) UseLua() bool {
	return false
}

func (c CatCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	placeholder := tgbotapi.NewPhotoShare(
		update.Message.Chat.ID,
		"https://static.thenounproject.com/png/101791-200.png",
	)
	placeholder.Caption = "ищу котека..."
	placeholderMessage, _ := api.Send(placeholder)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	media := &tgbotapi.BaseInputMedia{
		Type:    "",
		Media:   "",
		Caption: "котек найден!",
	}

	isGif := false
	if rand.Intn(10) <= 2 {
		media.Media = fmt.Sprintf("https://thecatapi.com/api/images/get?type=gif&%s", timestamp)
		media.Type = "animation"
		isGif = true
	} else {
		media.Media = fmt.Sprintf("https://thecatapi.com/api/images/get?type=jpg,png&%s", timestamp)
		media.Type = "photo"
	}

	editmsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    placeholderMessage.Chat.ID,
			MessageID: placeholderMessage.MessageID,
		},
		Media: media,
	}

	_, err := api.Send(editmsg)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to edit a message")
	}

	if isGif {
		metrics.Graphite().Send(graphite.NewMetric(
			fmt.Sprintf("stats.gif_cat.%s", update.Message.From.UserName), "1",
			time.Now().Unix(),
		))
	}
}
