package commands

import (
	"gopkg.in/telegram-bot-api.v4"
)

var BannedStickers map[string]bool

func init() {
	BannedStickers = map[string]bool{
		"CAADAgAD4QADtok9EpG19hLYFcFjAg":  true,
		"CAADAgADJAADtok9EplD2R-DctH5Ag":  true,
		"CAADAgAEAwAC4HlSB_6fe7DsL3ZdAg":  true,
		"CAADAgADgwMAAuB5UgcIctTnmytWOgI": true,
		"CAADAQADqwADmPThA6JMaCSxw5ePAg":  true,
	}
}

func Wlogys(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	stickerID := update.Message.ReplyToMessage.Sticker.FileID
	BannedStickers[stickerID] = true
}
