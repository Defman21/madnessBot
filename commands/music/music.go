package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/templates"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

type commandTemplate struct {
	Type  string
	Room  string
	Title string
	ID    string
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
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

	msg := templates.ExecuteTemplate("commands_music", commandTemplate{
		Type:  response.Data.CurrentSong.Type,
		Room:  room,
		Title: response.Data.CurrentSong.Name,
		ID:    response.Data.CurrentSong.ID,
	})
	helpers.SendMessage(api, update, msg, false, true)
}

func init() {
	commands.Register("music", &Command{})
}
