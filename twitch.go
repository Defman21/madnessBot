package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/golang-lru"
	"github.com/marpaia/graphite-golang"
	"github.com/nicklaw5/helix/v2"
	"io/ioutil"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/common/metrics"
	"madnessBot/config"
	"madnessBot/redis"
	"madnessBot/state/online"
	"madnessBot/templates"
	"net/http"
	"strconv"
	"time"
)

var notificationIds *lru.Cache

func init() {
	notificationIds, _ = lru.New(config.Config.NotificationsLRU)
}

type notificationTemplate struct {
	Login   string
	Title   string
	Viewers int
	Game    string
	UserID  string
}

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func twitchNotificationHandler(api *tgbotapi.BotAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := r.URL.Path[len(config.Config.Twitch.Webhook.Path):]
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		logger.Log.Debug().Bytes("body", body).Msg("body")

		if !helix.VerifyEventSubNotification(config.Config.Twitch.Webhook.Secret, r.Header,
			string(body)) {
			logger.Log.Error().Msg("invalid twitch signature on subscription")
			return
		} else {
			logger.Log.Debug().Msg("incoming twitch subscription")
		}

		var vals eventSubNotification
		err := json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to parse twitch subscription")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if vals.Challenge != "" {
			w.Write([]byte(vals.Challenge))
			return
		}

		subKey := fmt.Sprintf("%s:%s", channel, vals.Subscription.Type)
		subID := vals.Subscription.ID
		_, err = redis.Get().HSet(context.Background(), redis.HelixSubscriptionsKey, subKey, subID).Result()

		if err != nil {
			logger.Log.Error().Err(err).Str("channel", channel).Str("subscription-id",
				vals.Subscription.ID).Msg("Failed to save subscription")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var onlineEvent helix.EventSubStreamOnlineEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&onlineEvent)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to parse twitch `channel.online` event")
			return
		}

		// todo: `channel.offline` event support
		//message = templates.ExecuteTemplate("twitch_stream_ended", struct {
		//	Login string
		//}{Login: channel})
		//helpers.SendMessageChatID(api, config.Config.ChatID, message)
		//online.Add(channel, false)

		if _, exists := notificationIds.Get(r.Header.Get("Twitch-Eventsub-Message-Id")); exists {
			logger.Log.Info().Msg("Duplicate notification")
			return
		}

		notificationIds.Add(r.Header.Get("Twitch-Eventsub-Message-Id"), true)

		stream, err := config.Config.Twitch.Client().GetStreams(&helix.StreamsParams{
			UserIDs: []string{onlineEvent.BroadcasterUserID},
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to get stream")
			return
		}
		streamData := stream.Data.Streams
		if len(streamData) == 0 {
			return
		}

		message := templates.ExecuteTemplate(
			"twitch_stream_started",
			notificationTemplate{
				Login:   channel,
				Title:   streamData[0].Title,
				Viewers: streamData[0].ViewerCount,
				Game:    streamData[0].GameName,
				UserID:  onlineEvent.BroadcasterUserID,
			},
		)

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		url := "https://static-cdn.jtvnw.net/previews-ttv/live_user_" +
			channel + "-1280x720.jpg?t=" + timestamp
		helpers.SendPhotoChatID(api, config.Config.ChatID, url, message)

		metrics.Graphite().Send(graphite.NewMetric(
			fmt.Sprintf("stats.stream_push.%s", channel), "1",
			time.Now().Unix(),
		))

		online.Add(channel, true)
		return
	}
}
