package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/helpers"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/Defman21/madnessBot/common/types"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"os"
)

type Command struct{}
type Users map[string]string

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if !common.IsAdmin(update.Message.From) {
		helpers.SendMessage(api, update, "TriHard LULW", true)
		return
	}

	channel := update.Message.CommandArguments()

	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}

	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		common.Log.Warn().Err(err).Msg("Failed to read users.json")
		return
	}

	var users Users

	json.Unmarshal(bytes, &users)

	if userID, ok := users[channel]; ok {
		go func(channel string, userID string) {
			req := helpers.Request.Post("https://api.twitch.tv/helix/webhooks/hub").Query(
				types.TwitchHub{
					Callback:     fmt.Sprintf("%s%s", os.Getenv("TWITCH_URL"), channel),
					Mode:         "unsubscribe",
					LeaseSeconds: 864000,
					Topic:        fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID),
				},
			)
			oauth.AddHeadersUsing("twitch", req)
			_, _, errs := req.End()

			if errs != nil {
				common.Log.Error().Errs("errs", errs).Msg("Failed to send a request")
				return
			}

			common.Log.Info().
				Str("user", channel).Msg("Unsubscribed")

			delete(users, channel)
			jsonStr, _ := json.Marshal(users)
			err = ioutil.WriteFile("./data/users.json", []byte(jsonStr), 0644)
			if err == nil {
				common.Log.Info().Msg("Updated users.json")
				helpers.SendMessage(
					api,
					update,
					fmt.Sprintf("Unsubscribed from %s", channel),
					true,
				)
			} else {
				common.Log.Warn().Err(err).Msg("Couldn't write to users.json")
			}
		}(channel, userID)
	} else {
		common.Log.Warn().Str("channel", channel).Msg("Channel not found")
	}
}

func init() {
	commands.Register("unsubscribe", &Command{})
}
