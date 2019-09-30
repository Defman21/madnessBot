package main

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/metrics"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/Defman21/madnessBot/notifier"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/golang-lru"
	"github.com/marpaia/graphite-golang"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var notificationIds *lru.Cache

func init() {
	var cacheSize int
	cacheSizeEnv := os.Getenv("NOTIFICATIONS_LRU_CACHE")
	if cacheSizeEnv == "" {
		cacheSize = 10
	} else {
		cacheSize, _ = strconv.Atoi(cacheSizeEnv)
	}
	notificationIds, _ = lru.New(cacheSize)
}

func madnessTwitch(api *tgbotapi.BotAPI) http.HandlerFunc {
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
				helpers.SendMessageChatID(api, chatID, message)
			} else {

				if _, exists := notificationIds.Get(notification.Data[0].NotificationID); exists {
					common.Log.Info().Msg("Duplicate notification!")
					return
				}
				notificationIds.Add(notification.Data[0].NotificationID, true)

				type game struct {
					Name string
				}

				type gameResponse struct {
					Data []*game
				}

				var data gameResponse

				if notification.Data[0].Game != "0" {
					req := helpers.Request.Get("https://api.twitch.tv/helix/games").Query(
						struct {
							ID string
						}{
							ID: notification.Data[0].Game,
						},
					)
					oauth.AddHeadersUsing("twitch", req)
					_, _, errs := req.EndStruct(&data)

					if errs != nil {
						common.Log.Error().Errs("errs", errs).Msg("Request failed")
						return
					}
				} else {
					data = gameResponse{
						Data: []*game{{Name: "не указана"}},
					}
				}

				tpl := `%s завел подрубочку!
%s
Сморков: %d
Игра: %s
https://twitch.tv/%s
%s
`
				notifyUsers := notifier.Get().GenerateNotifyString(notification.Data[0].ID)
				message = fmt.Sprintf(tpl, name, notification.Data[0].Title,
					notification.Data[0].Viewers, data.Data[0].Name, name, notifyUsers)
				timestamp := strconv.FormatInt(time.Now().Unix(), 10)
				url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
					name + "-1280x720.jpg?" + timestamp
				helpers.SendPhotoChatID(api, chatID, url, message)

				metrics.Graphite().Send(graphite.NewMetric(
					fmt.Sprintf("stats.stream_push.%s", name), "1",
					time.Now().Unix(),
				))
			}
		}
	}
}
