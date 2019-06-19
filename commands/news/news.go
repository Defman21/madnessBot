package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"os"
	"time"

	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"gopkg.in/telegram-bot-api.v4"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	res, err := goreq.Request{
		Uri: "https://api.vk.com/method/wall.get",
		QueryString: struct {
			OwnerID     int `url:"owner_id"`
			Count       int
			Version     float64 `url:"v"`
			AccessToken string  `url:"access_token"`
		}{
			OwnerID:     -30138776,
			Count:       2,
			Version:     5.71,
			AccessToken: os.Getenv("VK_TOKEN"),
		},
	}.Do()

	if err != nil {
		common.Log.Error().Err(err).Msg("Request failed")
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
	res.Body.FromJsonTo(&data)

	if data.Response.Items[0].Pinned == 1 {
		data.Response.Items[0] = data.Response.Items[1]
	}

	url := fmt.Sprintf("https://vk.com/wall%d_%d", data.Response.Items[0].OwnerID, data.Response.Items[0].ID)

	text := fmt.Sprintf("%s\n%s\n%s", time.Unix(data.Response.Items[0].Date, 0).Format("02.01 15:04"), data.Response.Items[0].Text, url)

	for _, attachment := range data.Response.Items[0].Attachments {
		if attachment.Type != "photo" {
			continue
		} else {
			photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
			photo.FileID = attachment.Photo.URL
			photo.UseExisting = true
			photo.Caption = text
			bot.Send(photo)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}

func init() {
	commands.Register("news", &Command{})
}
