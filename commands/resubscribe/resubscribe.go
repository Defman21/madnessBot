package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/franela/goreq"
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
		return
	}
	bytes, err := ioutil.ReadFile("./data/users.json")
	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to read users.json")
		return
	}

	var users Users

	json.Unmarshal(bytes, &users)

	for channel, userID := range users {
		go func(channel string, userID string) {
			req := goreq.Request{
				Method: "POST",
				Uri:    "https://api.twitch.tv/helix/webhooks/hub",
				QueryString: struct {
					HubCallback     string `url:"hub.callback"`
					HubMode         string `url:"hub.mode"`
					HubLeaseSeconds int    `url:"hub.lease_seconds"`
					HubTopic        string `url:"hub.topic"`
				}{
					HubCallback:     fmt.Sprintf("%s%s", os.Getenv("TWITCH_URL"), channel),
					HubMode:         "subscribe",
					HubLeaseSeconds: 864000,
					HubTopic:        fmt.Sprintf("https://api.twitch.tv/helix/streams?user_id=%s", userID),
				},
			}
			oauth.AddHeadersUsing("twitch", &req)
			_, err := req.Do()

			if err != nil {
				common.Log.Error().Err(err).Msg("Failed to resubscribe")
			} else {
				common.Log.Info().Str("channel", channel).Msg("Subscribed")
			}
		}(channel, userID)
	}

	common.ResubscribeState.Save()
}

func init() {
	commands.Register("resubscribe", &Command{})
}