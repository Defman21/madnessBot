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

func Resubscribe(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		common.Log.Info("Prevented resubscribe")
		return
	}
	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		common.Log.Warn(err.Error())
		return
	}

	var users Users

	json.Unmarshal(bytes, &users)

	common.Log.WithFields(logrus.Fields{
		"users": users,
	}).Info("Users")

	for channel, userID := range users {
		common.Log.WithFields(logrus.Fields{
			"channel": channel,
			"user_id": userID,
		}).Info("Resubbing")
		go func(channel string, userID string) {
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
					HubTopic:        fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID),
				},
			}
			req.AddHeader("Client-ID", os.Getenv("TWITCH_TOKEN"))
			_, err := req.Do()

			if err != nil {
				common.Log.Warn(err.Error())
			} else {
				common.Log.WithFields(logrus.Fields{
					"user": channel,
				}).Info("Subscribed")
			}
		}(channel, userID)
	}
}
