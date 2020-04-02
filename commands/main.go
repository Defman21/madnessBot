package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marpaia/graphite-golang"
	"madnessBot/common/logger"
	"madnessBot/common/metrics"
	"time"
)

type Command interface {
	Run(api *tgbotapi.BotAPI, update *tgbotapi.Update)
	UseLua() bool
}

var commands = make(map[string]Command)

func Register(name string, command Command) {
	if _, exists := commands[name]; !exists {
		commands[name] = command
	}
}

func Run(name string, api *tgbotapi.BotAPI, update *tgbotapi.Update) bool {
	if command, exists := commands[name]; exists {
		if command.UseLua() {
			// TODO: Run {name}/{name}.lua
		}
		logger.Log.Info().Str("command", name).Msg("Called a command")
		metrics.Graphite().Send(graphite.NewMetric(
			fmt.Sprintf("stats.command.%s", name), "1",
			time.Now().Unix(),
		))
		go command.Run(api, update)
		return true
	} else {
		return false
	}
}
