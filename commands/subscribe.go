package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"os"
)

type Users map[string]string

func Subscribe(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	channel := update.Message.CommandArguments()
	if channel == "" {
		msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID,
			"AwADAgADwgADC6ZpS13yfdzm_pTzAg")
		bot.Send(msg)
		return
	}
	req := goreq.Request{
		Uri: "https://api.twitch.tv/helix/users",
		QueryString: struct {
			Login string
		}{
			Login: channel,
		},
	}
	req.AddHeader("Client-ID", os.Getenv("TWITCH_TOKEN"))
	res, err := req.Do()

	if err != nil {
		common.Log.Error().Err(err).Msg("Request failed")
		return
	} else {
		type User struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		var user User

		res.Body.FromJsonTo(&user)

		if len(user.Data) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такого пидора нет")
			bot.Send(msg)
			return
		}

		req := goreq.Request{
			Method: "POST",
			Uri:    "https://api.twitch.tv/helix/webhooks/hub",
			QueryString: struct {
				HubCallback     string `url:"hub.callback"`
				HubMode         string `url:"hub.mode"`
				HubLeaseSeconds int    `url:"hub.lease_seconds"`
				HubTopic        string `url:"hub.topic"`
			}{
				HubCallback:     fmt.Sprintf("%s%s", os.Getenv("TWITCH_URL"), channel),
				HubMode:         "subscribe",
				HubLeaseSeconds: 864000,
				HubTopic: fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s",
					user.Data[0].ID),
			},
		}
		req.AddHeader("Client-ID", os.Getenv("TWITCH_TOKEN"))
		_, err := req.Do()

		var users Users
		bytes, err := ioutil.ReadFile("./data/users.json")
		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to read users.json")
		}

		json.Unmarshal(bytes, &users)

		users[channel] = user.Data[0].ID
		bytes, err = json.Marshal(users)

		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to serialize users")
		} else {
			err := ioutil.WriteFile("./data/users.json", bytes, 0644)
			if err != nil {
				common.Log.Error().Err(err).Msg("Failed to write users.json")
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				fmt.Sprintf("Бот теперь аки маньяк будет преследовать %s "+
					"до конца своих дней.",
					channel))
			bot.Send(msg)
		}
	}
}
