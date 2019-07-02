package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"strconv"
	"strings"
	"time"

	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/franela/goreq"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	channel := update.Message.CommandArguments()

	if channel == "" {
		msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID,
			"AwADAgADwgADC6ZpS13yfdzm_pTzAg")
		api.Send(msg)
		return
	}

	channel = strings.ToLower(channel)

	placeholder := tgbotapi.NewPhotoShare(
		update.Message.Chat.ID,
		"https://static.thenounproject.com/png/101791-200.png",
	)
	placeholder.Caption = "ищу стримера..."
	placeholderMsg, err := api.Send(placeholder)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to send a placeholder message")
		return
	}

	req := goreq.Request{
		Uri: "https://api.twitch.tv/helix/streams",
		QueryString: struct {
			UserLogin string `url:"user_login"`
		}{
			UserLogin: channel,
		},
	}
	oauth.AddHeadersUsing("twitch", &req)
	res, err := req.Do()

	if err != nil {
		common.Log.Error().Err(err).Msg("Request failed")
		return
	}
	type TwitchResponse struct {
		Data []struct {
			Title   string `json:"title"`
			Viewers int64  `json:"viewer_count"`
			Game    string `json:"game_id"`
		} `json:"data"`
	}

	var data TwitchResponse
	res.Body.FromJsonTo(&data)

	editmsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    placeholderMsg.Chat.ID,
			MessageID: placeholderMsg.MessageID,
		},
	}

	if len(data.Data) != 0 {
		type game struct {
			Name string `json:"name"`
		}

		type gameResponse struct {
			Data []game `json:"data"`
		}

		var gdata gameResponse

		if data.Data[0].Game != "0" {
			req = goreq.Request{
				Uri: "https://api.twitch.tv/helix/games",
				QueryString: struct {
					ID string
				}{
					ID: data.Data[0].Game,
				},
			}
			oauth.AddHeadersUsing("twitch", &req)
			res, err = req.Do()

			if err != nil {
				common.Log.Error().Err(err).Msg("Request failed")
				return
			}

			res.Body.FromJsonTo(&gdata)
		} else {
			gdata = gameResponse{Data: []game{{Name: "не указана"}}}
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
			channel + "-1280x720.jpg?" + timestamp
		msg := fmt.Sprintf(`%s сейчас онлайн!
%s
Сморков: %d
Игра: %s
https://twitch.tv/%s`, channel, data.Data[0].Title, data.Data[0].Viewers, gdata.Data[0].Name, channel)

		editmsg.Media = tgbotapi.BaseInputMedia{
			Type:    "photo",
			Media:   url,
			Caption: msg,
		}
	} else {
		editmsg.Media = tgbotapi.BaseInputMedia{
			Type:    "photo",
			Media:   "https://i.redd.it/07onk217ojfz.png",
			Caption: fmt.Sprintf("%s ниче не стримит", channel),
		}
	}

	_, err = api.Send(editmsg)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to edit a message")
	}
}

func init() {
	commands.Register("info", &Command{})
}
