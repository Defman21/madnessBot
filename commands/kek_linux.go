package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"net/http"
	"os"
)

type KekCmd struct{}

func (c KekCmd) UseLua() bool {
	return false
}

func (c KekCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !helpers.PayCheck(api, update) {
		return
	}
	photos, err := api.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(update.Message.From.ID))
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to get user profile photo")
	} else {
		direction := update.Message.CommandArguments()
		zulul := photos.Photos[0]
		photo := zulul[len(zulul)-1]
		url, _ := api.GetFileDirectURL(photo.FileID)
		img, _ := os.Create("zulul.jpg")
		defer img.Close()

		resp, _ := http.Get(url)
		defer resp.Body.Close()
		_, _ = io.Copy(img, resp.Body)
		imagick.Initialize()
		defer imagick.Terminate()
		mw := imagick.NewMagickWand()
		defer mw.Destroy()
		mw.ReadImage("zulul.jpg")
		w := mw.GetImageWidth()
		h := mw.GetImageHeight()

		mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_DEACTIVATE)
		if direction == "right" {
			mw.CropImage(w/2, h, int(w/2), 0)
			mw.FlopImage()
		} else {
			mw.CropImage(w/2, h, 0, 0)
		}
		mwr := mw.Clone()
		defer mwr.Destroy()
		mwr.FlopImage()
		mw.AddImage(mwr)
		mw.SetFirstIterator()

		mwout := mw.AppendImages(false)
		defer mwout.Destroy()
		mwout.WriteImage("zulul-done.jpg")

		api.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "zulul-done.jpg"))
	}
}
