package commands

import (
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"gopkg.in/telegram-bot-api.v4"
	"os"
)

func News(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	res, err := goreq.Request{
		Uri: "https://api.vk.com/method/wall.get",
		QueryString: struct {
			OwnerID     int `url:"owner_id"`
			Count       int
			Version     float64 `url:"v"`
			AccessToken string  `url:"access_token"`
		}{
			OwnerID:     -30138776,
			Count:       1,
			Version:     5.71,
			AccessToken: os.Getenv("VK_TOKEN"),
		},
	}.Do()

	if err != nil {
		common.Log.Warn(err.Error())
		return
	} else {
		type VkResponse struct {
			Response struct {
				Items []struct {
					Text        string `json:"text"`
					Attachments []struct {
						Photo struct {
							URL string `json:"photo_604"`
						} `json:"photo"`
					} `json:"attachments"`
				} `json:"items"`
			} `json:"response"`
		}

		var data VkResponse
		res.Body.FromJsonTo(&data)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, data.Response.Items[0].Text)
		bot.Send(msg)
		photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
		photo.FileID = data.Response.Items[0].Attachments[0].Photo.URL
		photo.UseExisting = true
		bot.Send(photo)
	}
}
