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

	stream, errs := helpers.GetTwitchStreamByLogin(channel)

	if errs != nil {
		logger.Log.Error().Errs("errs", errs).Msg("Failed to get the stream")
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
		game, errs := helpers.GetTwitchGame(stream.Game)
		if errs != nil {
			logger.Log.Error().Errs("errs", errs).Msg("Failed to get the game")
			return
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		infoCommand.Online = true
		infoCommand.Title = stream.Title
		infoCommand.Viewers = stream.Viewers

		if game != nil {
			infoCommand.Game = game.Name
		} else {
			infoCommand.Game = "не указана"
		}

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
		logger.Log.Error().Err(err).Msg("Failed to edit a message")
	}
}
