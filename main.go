package main

import (
	cmds "github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var log = common.Log

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

	commands := map[string]func(*tgbotapi.BotAPI, *tgbotapi.Update){
		"up":          cmds.Up,
		"news":        cmds.News,
		"dnevnik":     cmds.Dnevnik,
		"birthday":    cmds.Birthday,
		"fuck":        cmds.Swap,
		"info":        cmds.Info,
		"subscribe":   cmds.Subscribe,
		"cat":         cmds.Cat,
		"music":       cmds.Music,
		"me":          cmds.Me,
		"resubscribe": cmds.Resubscribe,
		"unsubscribe": cmds.Unsubscribe,
		"wlogys":      cmds.Wlogys,
	}

	if err != nil {
		log.WithFields(logrus.Fields{
			"token": os.Getenv("BOT_TOKEN"),
		}).Fatal(err)
	}

	log.Printf("Account name: %s", bot.Self.UserName)

	_, err = bot.RemoveWebhook()

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("MADNESS_URL")))
	if err != nil {
		log.WithFields(logrus.Fields{
			"url": os.Getenv("MADNESS_URL"),
		}).Fatal(err.Error())
	}

	info, err := bot.GetWebhookInfo()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.WithFields(logrus.Fields{
			"webhook": info,
		}).Info("Webhook set")
	}

	updates := bot.ListenForWebhook(os.Getenv("MADNESS_HOOK"))

	http.HandleFunc(os.Getenv("TWITCH_HOOK"), madnessTwitch(bot))

	go http.ListenAndServe("0.0.0.0:9000", nil)

	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	sleepRegex, err := regexp.Compile(`\A(?:я|Я)\s+спать`)
	sadRegex, err := regexp.Compile(`\A(?:я|Я)\s+обидел(?:ась|ся)`)

	if err != nil {
		log.Fatal(err.Error())
	}

	for update := range updates {
		log.WithFields(logrus.Fields{
			"update": update,
		}).Debug("Update")

		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != chatID {
			continue
		}

		if sticker := update.Message.Sticker; sticker != nil {
			if update.Message.From.ID == 370779007 {
				if _, banned := cmds.BannedStickers[sticker.FileID]; banned {
					_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.MessageID,
					})
					if err != nil {
						log.Warn(err.Error())
					}
				}
			}
		}

		command, exists := commands[update.Message.Command()]
		if exists {
			go command(bot, &update)
		} else {
			if sleepRegex.MatchString(update.Message.Text) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Споки <3"))
				commands["cat"](bot, &update)
			} else if sadRegex.MatchString(update.Message.Text) {
				msg := tgbotapi.NewStickerShare(update.Message.Chat.ID,
					"CAADAgAD9wIAAlwCZQO1cgzUpY4T7wI")
				bot.Send(msg)
			}
		}
	}
}
