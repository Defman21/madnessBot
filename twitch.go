package main

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/metrics"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/Defman21/madnessBot/templates"
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

type notificationData struct {
	NotificationID string `json:"id"`
	ID             string `json:"user_id"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	Game           string `json:"game_id"`
	Viewers        int    `json:"viewer_count"`
}

type Notification struct {
	Data []notificationData `json:"data"`
}

type twitchGameData struct {
	Name string
}

type gameRequest struct {
	ID string
}

type gameResponse struct {
	Data []*twitchGameData
}

type notificationTemplate struct {
	Login   string
	Title   string
	Viewers int
	Game    string
	UserID  string
}

func twitchNotificationHandler(api *tgbotapi.BotAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(os.Getenv("TWITCH_HOOK")):]
		bytes, _ := ioutil.ReadAll(r.Body)

		challenge := r.FormValue("hub.challenge")
		if len(challenge) > 1 {
			common.Log.Info().Str("name", name).Str("challenge", challenge).Msg("Challenge")

			_, _ = w.Write([]byte(challenge))
		} else {
			var notification Notification
			_ = json.Unmarshal(bytes, &notification)

			common.Log.Info().Interface("notification", notification).Msg("Notification")

			chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
			var message string
			if len(notification.Data) == 0 {
				message = templates.ExecuteTemplate("twitch_stream_ended", struct {
					Login string
				}{Login: name})
				helpers.SendMessageChatID(api, chatID, message)
			} else {
				data := notification.Data[0]
				if _, exists := notificationIds.Get(data.NotificationID); exists {
					common.Log.Info().Msg("Duplicate notification!")
					return
				}
				notificationIds.Add(data.NotificationID, true)

				var game twitchGameData
				var gameResp gameResponse

				if data.Game != "0" {
					req := helpers.Request.Get("https://api.twitch.tv/helix/games").Query(gameRequest{
						ID: data.Game,
					})
					oauth.AddHeadersUsing("twitch", req)
					_, _, errs := req.EndStruct(&gameResp)

					if errs != nil {
						common.Log.Error().Errs("errs", errs).Msg("Request failed")
						return
					}
					game = *gameResp.Data[0]
				} else {
					game = twitchGameData{Name: "не указана"}
				}
				message = templates.ExecuteTemplate(
					"twitch_stream_started",
					notificationTemplate{
						Login:   name,
						Title:   data.Title,
						Viewers: data.Viewers,
						Game:    game.Name,
						UserID:  data.ID,
					},
				)

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
