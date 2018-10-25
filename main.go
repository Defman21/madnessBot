package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	cmds "github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

var log = common.Log

func main() {
	noWebhook := flag.Bool("nowebhook", false, "Don't use webhooks")
	flag.Parse()
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
		"quote":       cmds.Quote,
		"addquote":    cmds.AddQuote,
		"quotelist":   cmds.QuoteList,
		"reverse":     cmds.Reverse,
		"kek":         cmds.Kek,
		"s":           cmds.Sarcasm,
	}

	if err != nil {
		log.WithFields(logrus.Fields{
			"token": os.Getenv("BOT_TOKEN"),
		}).Fatal(err)
	}
	var updates tgbotapi.UpdatesChannel
	log.Printf("Account name: %s", bot.Self.UserName)
	if *noWebhook {
		_, _ = bot.RemoveWebhook()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 3

		updates, _ = bot.GetUpdatesChan(u)
	} else {
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

		updates = bot.ListenForWebhook(os.Getenv("MADNESS_HOOK"))
	}

	http.HandleFunc(os.Getenv("TWITCH_HOOK"), madnessTwitch(bot))

	go http.ListenAndServe("0.0.0.0:9000", nil)

	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	sleepRegex := regexp.MustCompile(`(?i)\Aя\s+спать`)
	sadRegex := regexp.MustCompile(`(?i)\Aя\s+обидел(?:ась|ся)`)
	wikiRegex := regexp.MustCompile(`(?i)^(?:что|кто) так(?:ое|ой|ая) ([^\?]+)`)

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

		log.WithFields(logrus.Fields{
			"uid":      update.Message.From.ID,
			"username": update.Message.From.UserName,
		}).Info("Message")

		//if sticker := update.Message.Sticker; sticker != nil {
		//	if update.Message.From.ID == 370779007 {
		//		if _, banned := cmds.BannedStickers[sticker.FileID]; banned {
		//			go func(chatid int64, msgid int) {
		//				_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		//					ChatID:    chatid,
		//					MessageID: msgid,
		//				})
		//				if err != nil {
		//					log.Warn(err.Error())
		//				}
		//			}(update.Message.Chat.ID, update.Message.MessageID)
		//		}
		//	}
		//}

		//if update.Message.From.ID == 370779007 {
		//	continue
		//}

		command, exists := commands[update.Message.Command()]
		if exists {
			go command(bot, &update)
		} else {
			if sleepRegex.MatchString(update.Message.Text) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Споки <3"))
				cmds.Cat(bot, &update)
			} else if sadRegex.MatchString(update.Message.Text) {
				msg := tgbotapi.NewStickerShare(update.Message.Chat.ID,
					"CAADAgAD9wIAAlwCZQO1cgzUpY4T7wI")
				bot.Send(msg)
			} else if lookup := wikiRegex.FindStringSubmatch(update.Message.Text); lookup != nil {
				req := goreq.Request{
					Uri:       "https://ru.wikipedia.org/w/api.php",
					UserAgent: "madnessBot (https://defman.me; me@defman.me) goreq",
					QueryString: struct {
						Action        string
						Titles        string
						Prop          string
						Explaintext   bool
						Exintro       bool
						Format        string
						Formatversion int
						Redirects     int
					}{
						Action:        "query",
						Titles:        lookup[1],
						Prop:          "extracts",
						Explaintext:   true,
						Exintro:       true,
						Format:        "json",
						Formatversion: 2,
						Redirects:     1,
					},
				}

				type response struct {
					Query struct {
						Pages []struct {
							Title   string `json:"title"`
							Extract string `json:"extract"`
						} `json:"pages"`
					} `json:"query"`
				}
				res, err := req.Do()
				if err != nil {
					log.WithFields(logrus.Fields{
						"err": err,
					}).Warn("Wikipedia lookup")
					continue
				}
				var data response
				err = res.Body.FromJsonTo(&data)
				if err != nil {
					log.WithFields(logrus.Fields{
						"err": err,
					}).Warn("json decode error")
				}
				if len(data.Query.Pages) != 0 && len(data.Query.Pages[0].Extract) != 0 {
					page := data.Query.Pages[0]
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("[%v](https://ru.wikipedia.org/wiki/%v) - %v\n", page.Title, page.Title, page.Extract))
					msg.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(msg)
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Википедия не знает forsenKek"))
				}
			}
		}
	}
}
