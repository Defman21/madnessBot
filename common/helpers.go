package common

import (
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/franela/goreq"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//GetTwitchUserIDByLogin get userID by Twitch login
func GetTwitchUserIDByLogin(login string) (string, bool) {
	req := goreq.Request{
		Uri: "https://api.twitch.tv/helix/users",
		QueryString: struct {
			Login string
		}{
			Login: login,
		},
	}
	oauth.AddHeadersUsing("twitch", &req)
	res, err := req.Do()

	if err != nil {
		Log.Error().Err(err).Msg("Request failed")
		return "", false
	}

	type User struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	var user User

	err = res.Body.FromJsonTo(&user)
	if err != nil {
		Log.Error().Err(err).Msg("Failed to unmarshal twitch response")
		return "", false
	}

	if len(user.Data) == 0 {
		return "", false
	}
	return user.Data[0].ID, true
}

func sendMessage(api *tgbotapi.BotAPI, message tgbotapi.Chattable) {
	_, err := api.Send(message)
	if err != nil {
		Log.Error().Err(err).Interface("msg", message).Msg("Failed to send a message")
	}
}

//SendMessage send a simple text message
func SendMessage(api *tgbotapi.BotAPI, chatID int64, text string, replyTo *int) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyTo != nil {
		msg.ReplyToMessageID = *replyTo
	}
	sendMessage(api, msg)
}

const dremoAVNDVoiceID = "AwADAgADwgADC6ZpS13yfdzm_pTzAg"

//SendInvalidArgumentsMessage send a voice message by DremoAVND
func SendInvalidArgumentsMessage(api *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewVoiceShare(chatID, dremoAVNDVoiceID)
	sendMessage(api, msg)
}
