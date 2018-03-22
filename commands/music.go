package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"gopkg.in/telegram-bot-api.v4"
	"time"
)

func Music(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	room := update.Message.CommandArguments()

	if room == "" {
		room = "melharucos"
	}

	res, err := goreq.Request{
		Uri:     fmt.Sprintf("https://api.dubtrack.fm/room/%s", room),
		Timeout: 3 * time.Second,
	}.Do()

	if err != nil {
		common.Log.Warn(err.Error())
		return
	} else {
		type Response struct {
			Data struct {
				CurrentSong struct {
					ID   string `json:"fkid"`
					Name string `json:"name"`
				} `json:"currentSong"`
			} `json:"data"`
		}

		var response Response
		res.Body.FromJsonTo(&response)

		if response.Data.CurrentSong.ID == "" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("В комнате %s тихо", room)))
			return
		}

		url := fmt.Sprintf("https://youtube.com/watch?v=%s", response.Data.CurrentSong.ID)

		msg := fmt.Sprintf("В комнате %s играет %s\n%s", room, response.Data.CurrentSong.Name, url)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
	}
}
