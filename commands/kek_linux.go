package commands

import (
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gopkg.in/telegram-bot-api.v4"
	"io"
	"net/http"
	"os"
)

// Kek lul
func Kek(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !payCheck(bot, update) {
		return
	}
	photos, err := bot.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(update.Message.From.ID))
	if err != nil {
		common.Log.Warn().Err(err).Msg("Failed to get user profile photo")
	} else {
		direction := update.Message.CommandArguments()
		zulul := photos.Photos[0]
		photo := zulul[len(zulul)-1]
		url, _ := bot.GetFileDirectURL(photo.FileID)
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

		bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "zulul-done.jpg"))
	}
}
