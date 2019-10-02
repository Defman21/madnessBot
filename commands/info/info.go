package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/templates"
	"strconv"
	"strings"
	"time"

	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

type commandTemplate struct {
	Login   string
	Title   string
	Viewers int
	Game    string
	Online  bool
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	channel := update.Message.CommandArguments()

	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
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

	type TwitchResponse struct {
		Data []struct {
			Title   string `json:"title"`
			Viewers int    `json:"viewer_count"`
			Game    string `json:"game_id"`
		} `json:"data"`
	}

	var data TwitchResponse
	req := helpers.Request.Get("https://api.twitch.tv/helix/streams").Query(
		struct {
			UserLogin string `json:"user_login"`
		}{
			UserLogin: channel,
		},
	)
	oauth.AddHeadersUsing("twitch", req)
	_, _, errs := req.EndStruct(&data)

	if errs != nil {
		common.Log.Error().Errs("errs", errs).Msg("Request failed")
		return
	}

	editmsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    placeholderMsg.Chat.ID,
			MessageID: placeholderMsg.MessageID,
		},
	}

	infoCommand := commandTemplate{
		Login:  channel,
		Online: false,
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
			req = helpers.Request.Get("https://api.twitch.tv/helix/games").Query(
				struct {
					ID string
				}{
					ID: data.Data[0].Game,
				},
			)
			oauth.AddHeadersUsing("twitch", req)
			_, _, errs = req.EndStruct(&gdata)

			if errs != nil {
				common.Log.Error().Errs("errs", errs).Msg("Request failed")
				return
			}
		} else {
			gdata = gameResponse{Data: []game{{Name: "не указана"}}}
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		infoCommand.Online = true
		infoCommand.Title = data.Data[0].Title
		infoCommand.Viewers = data.Data[0].Viewers
		infoCommand.Game = gdata.Data[0].Name

		url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
			channel + "-1280x720.jpg?" + timestamp

		msg := templates.ExecuteTemplate("commands_info", infoCommand)
		editmsg.Media = tgbotapi.BaseInputMedia{
			Type:      "photo",
			Media:     url,
			Caption:   msg,
			ParseMode: tgbotapi.ModeMarkdown,
		}
	} else {
		msg := templates.ExecuteTemplate("commands_info", infoCommand)
		editmsg.Media = tgbotapi.BaseInputMedia{
			Type:      "photo",
			Media:     "https://i.redd.it/07onk217ojfz.png",
			Caption:   msg,
			ParseMode: tgbotapi.ModeMarkdown,
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
