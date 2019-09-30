package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common/helpers"
	"time"

	"github.com/Defman21/madnessBot/common"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
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
		common.Log.Error().Errs("errs", errs).Msg("Request failed")
		return
	}

	if response.Data.CurrentSong.ID == "" {
		helpers.SendMessage(api, update, fmt.Sprintf("В комнате %s тихо", room), false)
		return
	}

	var url string
	if response.Data.CurrentSong.Type == "youtube" {
		url = fmt.Sprintf("https://youtube.com/watch?v=%s", response.Data.CurrentSong.ID)
	} else {
		url = ""
	}

	msg := fmt.Sprintf("В комнате %s играет %s\n%s", room, response.Data.CurrentSong.Name, url)
	helpers.SendMessage(api, update, msg, false)
}

func init() {
	commands.Register("music", &Command{})
}
