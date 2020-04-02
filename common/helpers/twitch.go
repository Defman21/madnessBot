package helpers

import (
	"madnessBot/common/logger"
	"madnessBot/common/oauth"
	"madnessBot/common/types"
	"madnessBot/config"
)

// GetTwitchUser get user by login
func GetTwitchUser(login string) (*types.TwitchUser, []error) {
	var response types.TwitchUserResponse

	req := Request.Get("https://api.twitch.tv/helix/users").Query(
		types.TwitchUserRequest{Login: login},
	)

	oauth.AddHeadersUsing("twitch", req)

	_, _, errs := req.EndStruct(&response)

	if len(response.Data) == 0 {
		return nil, errs
	}

	return &response.Data[0], errs
}

//GetTwitchUserIDByLogin get userID by Twitch login
func GetTwitchUserIDByLogin(login string) (string, bool) {
	user, errs := GetTwitchUser(login)

	if errs != nil {
		logger.Log.Error().Errs("errs", errs).Msg("Request failed")
		return "", false
	}

	if user != nil {
		return user.ID, true
	}

	return "", false
}

//SendTwitchHubMessage sends a message to the Twitch Hub
func SendTwitchHubMessage(channel string, mode string, topic string) []error {
	query := types.TwitchHub{
		Callback:     config.Config.Twitch.Webhook.GetURL(channel),
		Mode:         mode,
		LeaseSeconds: 864000,
		Topic:        topic,
	}

	req := Request.Post("https://api.twitch.tv/helix/webhooks/hub").Query(query)

	logger.Log.Debug().Interface("query", query).Msg("Twitch hub message")

	oauth.AddHeadersUsing("twitch", req)
	_, _, errs := req.End()

	return errs
}

func GetTwitchStreamByLogin(login string) (stream *types.TwitchStream, errs []error) {
	var response types.TwitchStreamResponse

	req := Request.Get("https://api.twitch.tv/helix/streams").Query(types.TwitchStreamRequest{
		UserLogin: login,
	})

	oauth.AddHeadersUsing("twitch", req)
	_, _, errs = req.EndStruct(&response)

	if len(response.Data) == 0 {
		return nil, errs
	}

	return &response.Data[0], errs
}

func GetTwitchGame(gameID string) (game *types.TwitchGame, errs []error) {
	var response types.TwitchGameResponse

	req := Request.Get("https://api.twitch.tv/helix/games").Query(
		types.TwitchGameRequest{ID: gameID},
	)

	oauth.AddHeadersUsing("twitch", req)
	_, _, errs = req.EndStruct(&response)

	if len(response.Data) == 0 {
		return nil, errs
	}

	return &response.Data[0], errs

}
