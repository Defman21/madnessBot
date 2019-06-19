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
	"gopkg.in/telegram-bot-api.v4"
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

	if len(data.Data) != 0 {
		type gameResponse struct {
			Data []struct {
				Name string `json:"name"`
			} `json:"data"`
		}

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

		var gdata gameResponse
		res.Body.FromJsonTo(&gdata)
		photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
			channel + "-1280x720.jpg?" + timestamp
		photo.FileID = url
		photo.UseExisting = true
		tpl := `%s сейчас онлайн!
%s
Сморков: %d
Игра: %s
https://twitch.tv/%s`
		photo.Caption = fmt.Sprintf(tpl, channel, data.Data[0].Title,
			data.Data[0].Viewers, gdata.Data[0].Name, channel)
		api.Send(photo)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Етот пидор ниче не стримит")
		api.Send(msg)
		sticker := tgbotapi.NewStickerShare(update.Message.Chat.ID,
			"CAADAgADIwAD43TSFjrD9SW8bXfjAg")
		api.Send(sticker)
	}
}

func init() {
	commands.Register("info", &Command{})
}
