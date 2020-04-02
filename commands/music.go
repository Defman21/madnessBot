package commands

import (
	"fmt"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/templates"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MusicCmd struct{}

func (c MusicCmd) UseLua() bool {
	return false
}

type musicCmdTemplate struct {
	Type  string
	Room  string
	Title string
	ID    string
}

func (c MusicCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	room := update.Message.CommandArguments()

	if room == "" {
		room = "melharucos"
	}

	type Response struct {
		Data struct {
			CurrentSong struct {
				ID   string `json:"fkid"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"currentSong"`
		} `json:"data"`
	}

	var response Response

	_, _, errs := helpers.Request.
		Get(fmt.Sprintf("https://api.dubtrack.fm/room/%s", room)).
		Timeout(3 * time.Second).
		EndStruct(&response)

	if errs != nil {
		logger.Log.Error().Errs("errs", errs).Msg("Request failed")
		return
	}

	if response.Data.CurrentSong.ID == "" {
		helpers.SendMessage(api, update, fmt.Sprintf("В комнате %s тихо", room), false, true)
		return
	}

	msg := templates.ExecuteTemplate("commands_music", musicCmdTemplate{
		Type:  response.Data.CurrentSong.Type,
		Room:  room,
		Title: response.Data.CurrentSong.Name,
		ID:    response.Data.CurrentSong.ID,
	})
	helpers.SendMessage(api, update, msg, false, true)
}
