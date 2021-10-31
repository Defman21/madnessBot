package commands

import (
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/templates"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InfoCmd struct{}

func (c InfoCmd) UseLua() bool {
	return false
}

type commandTemplate struct {
	Login   string
	Title   string
	Viewers int
	Game    string
	Online  bool
}

func (c InfoCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !helpers.PayCheck(api, update) {
		return
	}

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
		logger.Log.Error().Err(err).Msg("Failed to send a placeholder message")
		return
	}

	stream, err := helpers.GetTwitchStreamByLogin(channel)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get stream")
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

	if stream != nil {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		infoCommand.Online = true
		infoCommand.Title = stream.Title
		infoCommand.Viewers = stream.ViewerCount
		infoCommand.Game = stream.GameName

		url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
			channel + "-1280x720.jpg?" + timestamp

		msg := templates.ExecuteTemplate("commands_info", infoCommand)
		editmsg.Media = tgbotapi.BaseInputMedia{
			Type:      "photo",
			Media:     url,
			Caption:   msg,
			ParseMode: tgbotapi.ModeMarkdownV2,
		}
	} else {
		msg := templates.ExecuteTemplate("commands_info", infoCommand)
		editmsg.Media = tgbotapi.BaseInputMedia{
			Type:      "photo",
			Media:     "https://i.redd.it/07onk217ojfz.png",
			Caption:   msg,
			ParseMode: tgbotapi.ModeMarkdownV2,
		}
	}

	_, err = api.Send(editmsg)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to edit a message")
	}
}
