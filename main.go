package main

import (
	"flag"
	"fmt"
	"github.com/marpaia/graphite-golang"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	cmds "github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
	"github.com/joho/godotenv"
	"gopkg.in/telegram-bot-api.v4"
)

var log = common.Log

func main() {
	noWebhook := flag.Bool("nowebhook", false, "Don't use webhooks")
	useGraphite := flag.Bool("graphite", false, "Use graphite")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

	commands := map[string]func(*tgbotapi.BotAPI, *tgbotapi.Update){
		"up":          cmds.Up,
		"news":        cmds.News,
		"fuck":        cmds.Swap,
		"info":        cmds.Info,
		"subscribe":   cmds.Subscribe,
		"cat":         cmds.Cat,
		"music":       cmds.Music,
		"me":          cmds.Me,
		"resubscribe": cmds.Resubscribe,
		"unsubscribe": cmds.Unsubscribe,
		"subs":        cmds.Subscribers,
		"reverse":     cmds.Reverse,
		"kek":         cmds.Kek,
		"s":           cmds.Sarcasm,
		"version":     cmds.Version,
		"donate":      cmds.Donate,
	}

	if err != nil {
		log.Fatal().
			Err(err).
			Str("token", os.Getenv("BOT_TOKEN")).
			Msg("Failed to create a bot")
	}
	var updates tgbotapi.UpdatesChannel

	log.Info().Str("username", bot.Self.UserName).Msg("Connected")

	if *noWebhook {
		log.Debug().Msg("Long-polling")

		_, _ = bot.RemoveWebhook()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 3

		updates, _ = bot.GetUpdatesChan(u)

	} else {
		_, err = bot.RemoveWebhook()

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to remove a webhook")
		}

		_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("MADNESS_URL")))
		if err != nil {
			log.Fatal().
				Err(err).
				Str("url", os.Getenv("MADNESS_URL")).
				Msg("Failed to set a weebhok")
		}

		info, err := bot.GetWebhookInfo()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get webhook info")
		} else {
			log.Info().Interface("webhook", info).Msg("Webhook set")
		}

		updates = bot.ListenForWebhook(os.Getenv("MADNESS_HOOK"))
	}

	var Graphite *graphite.Graphite

	if *useGraphite {
		port, err := strconv.Atoi(os.Getenv("GRAPHITE_PORT"))

		if err != nil {
			log.Error().Err(err).Msg("Invalid GRAPHITE_PORT")
		}

		Graphite, err = graphite.NewGraphite(os.Getenv("GRAPHITE_HOST"), port)

		if err != nil {
			log.Error().Err(err).Msg("Failed to initialize graphite")
		}
	}

	http.HandleFunc(os.Getenv("TWITCH_HOOK"), madnessTwitch(bot, Graphite))

	go http.ListenAndServe("0.0.0.0:9000", nil)

	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	sleepRegex := regexp.MustCompile(`(?i)\Aя\s+спать`)
	sadRegex := regexp.MustCompile(`(?i)\Aя\s+обидел(?:ась|ся)`)
	wikiRegex := regexp.MustCompile(`(?i)^(?:что|кто) так(?:ое|ой|ая) ([^\?]+)`)

	go common.TwitchOauthState.Load()

	for update := range updates {
		if time.Now().Local().After(common.TwitchOauthState.ExpiresAt) {
			go common.TwitchOauthState.Refresh()
		}

		log.Debug().Interface("update", update).Msg("Update")

		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != chatID {
			continue
		}

		commandName := update.Message.Command()

		command, exists := commands[commandName]
		if exists {
			common.Log.Info().Str("command", commandName).Msg("Called a command")
			if *useGraphite {
				metric := graphite.NewMetric(
					fmt.Sprintf("stats.command.%s", commandName), "1",
					time.Now().Unix(),
				)
				err = Graphite.SendMetric(metric)
				if err != nil {
					log.Error().Err(err).Msg("Failed to send metric")
				}

			}
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
					log.Warn().Err(err).Msg("Wikipedia lookup")
					continue
				}
				var data response
				err = res.Body.FromJsonTo(&data)
				if err != nil {
					log.Warn().Err(err).Msg("Failed to decode JSON")
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
