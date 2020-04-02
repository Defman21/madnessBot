package main

import (
	"madnessBot/commands"
	"madnessBot/common/helpers"
	"madnessBot/common/logger"
	"madnessBot/common/oauth"
	"madnessBot/common/oauth/twitch"
	"madnessBot/config"
	"madnessBot/integrations/boosty"
	"madnessBot/integrations/wiki"
	"madnessBot/redis"
	"madnessBot/state/resubscribe"
	_ "madnessBot/templates"
	"net/http"
	"regexp"
	"time"

	"madnessBot/common/metrics"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	go config.Init()
	<-config.Initialized
	logger.SetLogLevel(config.Config.LogLevel)
	logger.Log.Info().Interface("config", config.Config).Msg("Initialized config")
}

const sadCatStickerID = "CAADAgAD9wIAAlwCZQO1cgzUpY4T7wI"

var log = &logger.Log

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Config.Token)

	redis.Init()
	oauth.Register("twitch", twitch.Instance)

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

	cmds := map[string]commands.Command{
		"cat":         commands.CatCmd{},
		"donate":      commands.DonateCmd{},
		"info":        commands.InfoCmd{},
		"kek":         commands.KekCmd{},
		"me":          commands.MeCmd{},
		"music":       commands.MusicCmd{},
		"news":        commands.NewsCmd{},
		"notify_me":   commands.NotifyMeCmd{},
		"online":      commands.OnlineCmd{},
		"resolve":     commands.ResolveCmd{},
		"r":           commands.ResolveCmd{},
		"resubscribe": commands.ResubscribeCmd{},
		"reverse":     commands.ReverseCmd{},
		"sarcasm":     commands.SarcasmCmd{},
		"subscribe":   commands.SubscribeCmd{},
		"subscribers": commands.SubscribersCmd{},
		"swap":        commands.SwapCmd{},
		"fuck":        commands.SwapCmd{},
		"unnotify_me": commands.UnnotifyMeCmd{},
		"unsubscribe": commands.UnsubscribeCmd{},
		"up":          commands.UpCmd{},
		"version":     commands.VersionCmd{},
	}

	for name, handler := range cmds {
		commands.Register(name, handler)
		log.Info().Str("command", name).Msg("Registered command")
	}

	sleepRegex := regexp.MustCompile(`(?i)\Aя\s+спать`)
	sadRegex := regexp.MustCompile(`(?i)\Aя\s+обидел(?:ась|ся)`)
	wikiRegex := regexp.MustCompile(`(?i)^(?:что|кто) так(?:ое|ой|ая) ([^?]+)`)

	for update := range updates {
		oauth.RefreshExpired()

		nextResubscribeCall := resubscribe.GetState()
		if nextResubscribeCall != nil && time.Now().Local().After(*nextResubscribeCall) {
			go commands.Run("resubscribe", bot, &update)
		}

		log.Debug().Interface("update", update).Msg("Update")

		chatID := update.Message.Chat.ID

		if update.Message == nil {
			continue
		}

		if chatID == config.Config.BoostyChatID {
			boosty.HandleUpdate(bot, &update)
			continue
		}

		if chatID != config.Config.ChatID && chatID != config.Config.BoostyChatID {
			continue
		}

		commandName := update.Message.Command()
		if ran := commands.Run(commandName, bot, &update); !ran {
			if sleepRegex.MatchString(update.Message.Text) {
				helpers.SendMessage(bot, &update, "Споки <3", true, false)
				go commands.Run("cat", bot, &update)
			} else if sadRegex.MatchString(update.Message.Text) {
				helpers.SendSticker(bot, &update, sadCatStickerID, false)
			} else if lookup := wikiRegex.FindStringSubmatch(update.Message.Text); lookup != nil {
				wiki.HandleUpdate(bot, &update, lookup)
			}
		}
	}
}
