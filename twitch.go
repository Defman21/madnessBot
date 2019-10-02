package main

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/metrics"
	"github.com/Defman21/madnessBot/common/types"
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
			return
		}

		var notificationRequest types.TwitchWebHookNotificationRequest
		_ = json.Unmarshal(bytes, &notificationRequest)

		common.Log.Info().Interface("notificationRequest", notificationRequest).Msg("Notification")

		chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)

		var message string

		if len(notificationRequest.Data) == 0 {
			message = templates.ExecuteTemplate("twitch_stream_ended", struct {
				Login string
			}{Login: name})
			helpers.SendMessageChatID(api, chatID, message)
			return
		}

		notification := notificationRequest.Data[0]

		if _, exists := notificationIds.Get(notification.ID); exists {
			common.Log.Info().Msg("Duplicate notificationRequest!")
			return
		}

		notificationIds.Add(notification.ID, true)

		game, errs := helpers.GetTwitchGame(notification.Game)
		if errs != nil {
			common.Log.Error().Errs("errs", errs).Msg("Failed to get the game")
		}

		if game == nil {
			game = &types.TwitchGame{Name: "не указана"}
		}

		message = templates.ExecuteTemplate(
			"twitch_stream_started",
			notificationTemplate{
				Login:   name,
				Title:   notification.Title,
				Viewers: notification.Viewers,
				Game:    game.Name,
				UserID:  notification.UserID,
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
