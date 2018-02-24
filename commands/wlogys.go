package commands

import (
	"github.com/Defman21/madnessBot/common"
	"github.com/sirupsen/logrus"
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
	defer func() {
		common.Log.Warn(recover())
	}()
	stickerID := update.Message.ReplyToMessage.Sticker.FileID
	BannedStickers[stickerID] = true
	common.Log.WithFields(logrus.Fields{
		"stickerID": stickerID,
	}).Info("Banned sticker for wlogys")
}
