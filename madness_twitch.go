package main

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"github.com/marpaia/graphite-golang"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var notificationIds map[string]bool

func init() {
	notificationIds = make(map[string]bool)
}

func madnessTwitch(bot *tgbotapi.BotAPI, graphiteSrv *graphite.Graphite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(os.Getenv("TWITCH_HOOK")):]
		type Notification struct {
			Data []struct {
				NotificationID string `json:"id"`
				ID             string `json:"user_id"`
				Title          string `json:"title"`
				Type           string `json:"type"`
				Game           string `json:"game_id"`
				Viewers        int    `json:"viewer_count"`
			} `json:"data"`
		}
		bytes, _ := ioutil.ReadAll(r.Body)

		challenge := r.FormValue("hub.challenge")
		if len(challenge) > 1 {
			common.Log.Info().Str("name", name).Str("challenge", challenge).Msg("Challenge")

			w.Write([]byte(challenge))
		} else {
			var notification Notification
			json.Unmarshal(bytes, &notification)

			common.Log.Info().Interface("notification", notification).Msg("Notification")

			chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
			var message string
			if len(notification.Data) == 0 {
				message = fmt.Sprintf("%s закончил стрим, потому что нахуй никому "+
					"не сдался.\nhttps://twitch.tv/%s", name, name)
				msg := tgbotapi.NewMessage(chatID, message)
				bot.Send(msg)
			} else {

				if _, dup := notificationIds[notification.Data[0].NotificationID]; dup {
					common.Log.Info().Msg("Duplicate notification!")
					return
				}
				notificationIds[notification.Data[0].NotificationID] = true
				req := goreq.Request{
					Uri: "https://api.twitch.tv/helix/games",
					QueryString: struct {
						ID string
					}{
						ID: notification.Data[0].Game,
					},
				}
				common.TwitchOauthState.AddHeaders(&req)
				res, err := req.Do()

				if err != nil {
					common.Log.Error().Err(err).Msg("Request failed")
					return
				}

				type gameResponse struct {
					Data []struct {
						Name string
					}
				}

				var data gameResponse
				res.Body.FromJsonTo(&data)

				tpl := `%s завел подрубочку!
%s
Сморков: %d
Игра: %s
https://twitch.tv/%s
`
				message = fmt.Sprintf(tpl, name, notification.Data[0].Title,
					notification.Data[0].Viewers, data.Data[0].Name, name)
				photo := tgbotapi.NewPhotoUpload(chatID, nil)
				timestamp := strconv.FormatInt(time.Now().Unix(), 10)
				url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
					name + "-1280x720.jpg?" + timestamp
				photo.FileID = url
				photo.UseExisting = true
				photo.Caption = message
				bot.Send(photo)
				metric := graphite.NewMetric(
					fmt.Sprintf("stats.stream_push.%s", name), "1",
					time.Now().Unix(),
				)

				if err := graphiteSrv.SendMetric(metric); err != nil {
					log.Error().Err(err).Msg("Failed to send metric")
				}
			}
		}
	}
}
