package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
)

type UpCmd struct{}

func (c UpCmd) UseLua() bool {
	return true
}

const upStickerFileID = "CAADAgADPgAD43TSFv8rTPYvm_MJAg"

func (c UpCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	helpers.SendSticker(api, update, upStickerFileID, true)
}
