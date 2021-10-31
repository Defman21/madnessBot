package helpers

import (
	"context"
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"madnessBot/common/logger"
	"madnessBot/config"
	"madnessBot/redis"
)

// GetTwitchUser get user by login
func GetTwitchUser(login string) (*helix.User, error) {
	resp, err := config.Config.Twitch.Client().GetUsers(&helix.UsersParams{
		Logins: []string{login},
	})

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get twitch user")
		return nil, err
	}

	if len(resp.Data.Users) == 0 {
		return nil, nil
	}

	return &resp.Data.Users[0], nil
}

//GetTwitchUserIDByLogin get userID by Twitch login
func GetTwitchUserIDByLogin(login string) (string, bool) {
	user, err := GetTwitchUser(login)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Request failed")
		return "", false
	}

	if user == nil {
		return "", false
	}

	return user.ID, true
}

//SendEventSubMessage sends a message to the Twitch Hub
func SendEventSubMessage(channel, eventType string) error {
	broadcasterID, success := GetTwitchUserIDByLogin(channel)
	if !success {
		logger.Log.Warn().Str("channel", channel).Msg("Channel not found")
		return nil
	}

	_, err := config.Config.Twitch.Client().CreateEventSubSubscription(&helix.
		EventSubSubscription{
		Type:    eventType,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: broadcasterID,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: config.Config.Twitch.Webhook.GetURL(channel),
			Secret:   config.Config.Twitch.Webhook.Secret,
		},
	})

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to send EventSub request")
		return err
	}

	return nil
}

//UnsubscribeFromEventSub unsubscribes from event
func UnsubscribeFromEventSub(channel, eventType string) error {
	subKey := fmt.Sprintf("%s:%s", channel, eventType)
	subId, err := redis.Get().HGet(context.Background(), redis.HelixSubscriptionsKey, subKey).Result()
	if err != nil {
		return err
	}

	_, err = config.Config.Twitch.Client().RemoveEventSubSubscription(subId)
	if err != nil {
		return err
	}

	n, err := redis.Get().HDel(context.Background(), redis.HelixSubscriptionsKey, subKey).Result()

	if err != nil {
		return err
	}

	if n != 1 {
		logger.Log.Warn().Int64("affected", n).Msg("too much affected redis keys")
	}

	return nil
}

func GetTwitchStreamByLogin(login string) (stream *helix.Stream, err error) {
	streams, err := config.Config.Twitch.Client().GetStreams(&helix.StreamsParams{
		UserLogins: []string{login},
	})

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get the stream")
		return nil, err
	}

	if len(streams.Data.Streams) == 0 {
		return nil, nil
	}

	return &streams.Data.Streams[0], err
}
