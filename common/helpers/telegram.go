package helpers

import (
	"github.com/Defman21/madnessBot/common/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMessage(api *tgbotapi.BotAPI, message tgbotapi.Chattable) {
	_, err := api.Send(message)
	if err != nil {
		logger.Log.Error().Err(err).Interface("msg", message).Msg("Failed to send a message")
	}
}

//SendMessage send a simple text message
func SendMessage(api *tgbotapi.BotAPI, update *tgbotapi.Update, text string, isReply bool) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if isReply {
		msg.ReplyToMessageID = update.Message.MessageID
	}
	msg.ParseMode = tgbotapi.ModeMarkdown
	sendMessage(api, msg)
}

//SendMessageChatID sends a message by chat id
func SendMessageChatID(api *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	sendMessage(api, msg)
}

//SendPhoto sends a photo with caption
func SendPhoto(api *tgbotapi.BotAPI, update *tgbotapi.Update, photoURL string, caption string, isReply bool) {
	photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
	photo.FileID = photoURL
	photo.UseExisting = true
	photo.Caption = caption
	if isReply {
		photo.ReplyToMessageID = update.Message.MessageID
	}
	photo.ParseMode = tgbotapi.ModeMarkdown
	sendMessage(api, photo)
}

//SendPhotoChatID sends a photo with caption by chat id
func SendPhotoChatID(api *tgbotapi.BotAPI, chatID int64, photoURL string, caption string) {
	photo := tgbotapi.NewPhotoUpload(chatID, nil)
	photo.FileID = photoURL
	photo.UseExisting = true
	photo.Caption = caption
	photo.ParseMode = tgbotapi.ModeMarkdown
	sendMessage(api, photo)
}

const dremoAVNDVoiceID = "AwADAgADwgADC6ZpS13yfdzm_pTzAg"

//SendInvalidArgumentsMessage send a voice message by DremoAVND
func SendInvalidArgumentsMessage(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewVoiceShare(update.Message.Chat.ID, dremoAVNDVoiceID)
	sendMessage(api, msg)
}
