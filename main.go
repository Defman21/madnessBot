package main

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	_ "github.com/Defman21/madnessBot/commands/cat"
	_ "github.com/Defman21/madnessBot/commands/donate"
	_ "github.com/Defman21/madnessBot/commands/info"
	_ "github.com/Defman21/madnessBot/commands/kek"
	_ "github.com/Defman21/madnessBot/commands/me"
	_ "github.com/Defman21/madnessBot/commands/music"
	_ "github.com/Defman21/madnessBot/commands/news"
	_ "github.com/Defman21/madnessBot/commands/notify_me"
	_ "github.com/Defman21/madnessBot/commands/online"
	_ "github.com/Defman21/madnessBot/commands/resolve"
	_ "github.com/Defman21/madnessBot/commands/resubscribe"
	_ "github.com/Defman21/madnessBot/commands/reverse"
	_ "github.com/Defman21/madnessBot/commands/sarcasm"
	_ "github.com/Defman21/madnessBot/commands/subscribe"
	_ "github.com/Defman21/madnessBot/commands/subscribers"
	_ "github.com/Defman21/madnessBot/commands/swap"
	_ "github.com/Defman21/madnessBot/commands/unnotify_me"
	_ "github.com/Defman21/madnessBot/commands/unsubscribe"
	_ "github.com/Defman21/madnessBot/commands/up"
	_ "github.com/Defman21/madnessBot/commands/version"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/Defman21/madnessBot/config"
	_ "github.com/Defman21/madnessBot/templates"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/metrics"

	_ "github.com/Defman21/madnessBot/common/oauth/twitch"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	if initialized := config.Init(); !initialized {
		os.Exit(1)
	}
	logger.SetLogLevel(config.Config.LogLevel)
	logger.Log.Info().Interface("cfg", config.Config).Msg("Initialized config")
}

var log = &logger.Log

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Config.Token)

	if err != nil {
		log.Fatal().
			Err(err).
			Str("token", config.Config.Token).
			Msg("Failed to create a bot")
	}
	var updates tgbotapi.UpdatesChannel

	log.Info().Str("username", bot.Self.UserName).Msg("Connected")

	if !config.Config.Webhook.Enabled() {
		log.Info().Str("method", "long-polling").Msg("Initialized Telegram API")
		_, _ = bot.Request(tgbotapi.RemoveWebhookConfig{})

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 3

		updates = bot.GetUpdatesChan(u)
	} else {
		log.Info().Str("method", "webhook").Msg("Initialized Telegram API")
		_, err = bot.Request(tgbotapi.RemoveWebhookConfig{})

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to remove a webhook")
		}

		_, err = bot.Request(tgbotapi.NewWebhook(config.Config.Webhook.GetURL()))
		if err != nil {
			log.Fatal().
				Err(err).
				Str("url", config.Config.Webhook.GetURL()).
				Msg("Failed to set a weebhok")
		}

		info, err := bot.GetWebhookInfo()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get webhook info")
		} else {
			log.Info().Interface("webhook", info).Msg("Webhook set")
		}

		updates = bot.ListenForWebhook(config.Config.Webhook.Path)
	}

	if config.Config.Graphite != nil {
		if config.Config.Graphite.Enabled {
			metrics.Init()
		}
	} else {
		log.Info().Msg("Graphite integration is disabled")
	}

	if config.Config.Twitch.Webhook.Enabled() {
		http.HandleFunc(config.Config.Twitch.Webhook.Path, twitchNotificationHandler(bot))
	} else {
		log.Info().Msg("Twitch integration is disabled")
	}

	go http.ListenAndServe(config.Config.Server.GetBindAddress(), nil)

	sleepRegex := regexp.MustCompile(`(?i)\Aя\s+спать`)
	sadRegex := regexp.MustCompile(`(?i)\Aя\s+обидел(?:ась|ся)`)
	wikiRegex := regexp.MustCompile(`(?i)^(?:что|кто) так(?:ое|ой|ая) ([^?]+)`)
	postyroniumRegex := regexp.MustCompile(`(?i)постирони(?:я|ю|и|й)`)

	go common.ResubscribeState.Load()

	for update := range updates {
		oauth.RefreshExpired()

		if time.Now().Local().After(common.ResubscribeState.ExpiresAt) {
			go commands.Run("resubscribe", bot, &update)
		}

		log.Debug().Interface("update", update).Msg("Update")

		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != config.Config.ChatID {
			continue
		}

		commandName := update.Message.Command()
		if ran := commands.Run(commandName, bot, &update); !ran {
			if postyroniumRegex.MatchString(update.Message.Text) {
				helpers.SendMessage(bot, &update, "постирай трусы", true)
			} else if sleepRegex.MatchString(update.Message.Text) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Споки <3"))
				commands.Run("cat", bot, &update)
			} else if sadRegex.MatchString(update.Message.Text) {
				msg := tgbotapi.NewStickerShare(update.Message.Chat.ID,
					"CAADAgAD9wIAAlwCZQO1cgzUpY4T7wI")
				bot.Send(msg)
			} else if lookup := wikiRegex.FindStringSubmatch(update.Message.Text); lookup != nil {
				type response struct {
					Query struct {
						Pages []struct {
							Title   string `json:"title"`
							Extract string `json:"extract"`
						} `json:"pages"`
					} `json:"query"`
				}
				var data response

				req := helpers.Request.Get("https://ru.wikipedia.org/w/api.php").
					Set("User-Agent", "madnessBot (https://defman.me; me@defman.me) gorequest").
					Query(
						struct {
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
					)

				_, _, errs := req.EndStruct(&data)
				if errs != nil {
					log.Warn().Errs("errs", errs).Msg("Wikipedia lookup")
					continue
				}
				if len(data.Query.Pages) != 0 && len(data.Query.Pages[0].Extract) != 0 {
					page := data.Query.Pages[0]
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("[%v](https://ru.wikipedia.org/wiki/%v) - %v\n", page.Title, page.Title, page.Extract))
					msg.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(msg)
				} else {
					helpers.SendMessage(bot, &update, "Википедия не знает forsenKek", true)
				}
			}
		}
	}
}
