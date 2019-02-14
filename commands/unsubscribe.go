package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"os"
)

func Unsubscribe(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	channel := update.Message.CommandArguments()

	if channel == "" {
		msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID,
			"AwADAgADwgADC6ZpS13yfdzm_pTzAg")
		bot.Send(msg)
		return
	}

	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		common.Log.Warn(err.Error())
		return
	}

	var users Users
	var newUsers Users
	findUser := false

	json.Unmarshal(bytes, &users)

	for _, user := range users.List {
		if user[0] != channel {
			newUsers.List = append(newUsers.List, user)
			continue
		} else {
			findUser = true
		}

		go func(user []string) {
			req := goreq.Request{
				Method: "POST",
				Uri:    "https://api.twitch.tv/helix/webhooks/hub",
				QueryString: struct {
					HubCallback     string `url:"hub.callback"`
					HubMode         string `url:"hub.mode"`
					HubLeaseSeconds int    `url:"hub.lease_seconds"`
					HubTopic        string `url:"hub.topic"`
				}{
					HubCallback:     fmt.Sprintf("%s%s", os.Getenv("TWITCH_URL"), user[0]),
					HubMode:         "unsubscribe",
					HubLeaseSeconds: 864000,
					HubTopic:        fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", user[1]),
				},
			}
			req.AddHeader("Client-ID", os.Getenv("TWITCH_TOKEN"))
			_, err := req.Do()

			if err != nil {
				common.Log.Warn(err.Error())
			} else {
				common.Log.WithFields(logrus.Fields{
					"user":    user[0],
					"context": "commands/unsubscribe",
				}).Info("Unsubscribed")

				jsonStr, _ := json.Marshal(newUsers)

				err = ioutil.WriteFile("./data/users.json", []byte(jsonStr), 0644)
				if err == nil {
					common.Log.Info("Updated users.json")
					bot.Send(
						tgbotapi.NewMessage(
							update.Message.Chat.ID,
							fmt.Sprintf("Unsubscribed from %s", channel),
						),
					)
				} else {
					common.Log.Warn("Couldn't write to users.json")
				}
			}
		}(user)
	}

	if findUser == false {
		common.Log.WithFields(logrus.Fields{
			"user": channel,
		}).Warn("Channel not found")
	}
}
