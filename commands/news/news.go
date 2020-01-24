package commands

import (
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/config"
	"github.com/Defman21/madnessBot/templates"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

type commandTemplate struct {
	Text    string
	Time    string
	OwnerID int64
	ID      int64
}

var nameToOwnerID = map[string]int64{}

func generateNameToOwnerMap() {
	for name, ownerID := range config.Config.News.Sources {
		nameToOwnerID[name] = ownerID
	}
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if config.Config.News == nil || !config.Config.News.Enabled {
		logger.Log.Info().Msg("VK News integration is disabled")
		return
	}

	generateNameToOwnerMap()
	name := update.Message.CommandArguments()

	if name == "" {
		name = "melharucos"
	}

	ownerID, exists := nameToOwnerID[name]

	if !exists {
		logger.Log.Warn().Interface("map", nameToOwnerID).Msg("Name does not exist in the map")
		return
	}

	type VkResponse struct {
		Response struct {
			Items []struct {
				Text        string `json:"text"`
				OwnerID     int64  `json:"owner_id"`
				ID          int64  `json:"id"`
				Date        int64  `json:"date"`
				Pinned      int64  `json:"is_pinned"`
				Attachments []struct {
					Type  string `json:"type"`
					Photo struct {
						URL string `json:"photo_604"`
					} `json:"photo"`
				} `json:"attachments"`
			} `json:"items"`
		} `json:"response"`
	}

	var data VkResponse

	_, _, errs := helpers.Request.Get("https://api.vk.com/method/wall.get").Query(
		struct {
			OwnerID     int64 `json:"owner_id"`
			Count       int
			Version     float64 `json:"v"`
			AccessToken string  `json:"access_token"`
		}{
			OwnerID:     ownerID,
			Count:       2,
			Version:     5.71,
			AccessToken: config.Config.News.Token,
		},
	).EndStruct(&data)

	if errs != nil {
		logger.Log.Error().Errs("errs", errs).Msg("Request failed")
		return
	}

	if data.Response.Items[0].Pinned == 1 {
		data.Response.Items[0] = data.Response.Items[1]
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	postTime := time.Unix(data.Response.Items[0].Date, 0).In(loc).Format("02.01 15:04")

	text := templates.ExecuteTemplate("commands_news", commandTemplate{
		Text:    data.Response.Items[0].Text,
		Time:    postTime,
		OwnerID: data.Response.Items[0].OwnerID,
		ID:      data.Response.Items[0].ID,
	})

	for _, attachment := range data.Response.Items[0].Attachments {
		if attachment.Type != "photo" {
			continue
		} else {
			if len(text) >= 1024 {
				helpers.SendPhoto(api, update, attachment.Photo.URL, "", false)
				helpers.SendMessage(api, update, text, false, true)
				return
			}
			helpers.SendPhoto(api, update, attachment.Photo.URL, text, false)
			return
		}
	}

	helpers.SendMessage(api, update, text, false, true)
}

func init() {
	commands.Register("news", &Command{})
}
