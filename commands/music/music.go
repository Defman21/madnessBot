package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"time"

	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"gopkg.in/telegram-bot-api.v4"
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

	res, err := goreq.Request{
		Uri:     fmt.Sprintf("https://api.dubtrack.fm/room/%s", room),
		Timeout: 3 * time.Second,
	}.Do()

	if err != nil {
		common.Log.Error().Err(err).Msg("Request failed")
		return
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
	res.Body.FromJsonTo(&response)

	if response.Data.CurrentSong.ID == "" {
		api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("В комнате %s тихо", room)))
		return
	}
	var url string
	if response.Data.CurrentSong.Type == "youtube" {
		url = fmt.Sprintf("https://youtube.com/watch?v=%s", response.Data.CurrentSong.ID)
	} else {
		url = ""
	}

	msg := fmt.Sprintf("В комнате %s играет %s\n%s", room, response.Data.CurrentSong.Name, url)
	api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
}

func init() {
	commands.Register("music", &Command{})
}
