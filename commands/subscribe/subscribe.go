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
	channel := update.Message.CommandArguments()
	if channel == "" {
		helpers.SendInvalidArgumentsMessage(api, update)
		return
	}
	userID, found := helpers.GetTwitchUserIDByLogin(channel)
	if found {
		req := helpers.Request.Post("https://api.twitch.tv/helix/webhooks/hub").Query(
			types.TwitchHub{
				Callback:     fmt.Sprintf("%s%s", os.Getenv("TWITCH_URL"), channel),
				Mode:         "subscribe",
				LeaseSeconds: 864000,
				Topic:        fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID),
			},
		)
		oauth.AddHeadersUsing("twitch", req)
		_, _, errs := req.End()

		if errs != nil {
			common.Log.Error().Errs("errs", errs).Msg("Failed to subscribe")
			return
		}

		var users Users
		bytes, err := ioutil.ReadFile("./data/users.json")
		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to read users.json")
		}

		json.Unmarshal(bytes, &users)

		users[channel] = userID
		bytes, err = json.Marshal(users)

		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to serialize users")
		} else {
			err := ioutil.WriteFile("./data/users.json", bytes, 0644)
			if err != nil {
				common.Log.Error().Err(err).Msg("Failed to write users.json")
				return
			}
			helpers.SendMessage(
				api,
				update,
				fmt.Sprintf("Бот теперь аки маньяк будет преследовать %s "+
					"до конца своих дней.", channel),
				true,
			)
		}
	}
}

func init() {
	commands.Register("subscribe", &Command{})
}
