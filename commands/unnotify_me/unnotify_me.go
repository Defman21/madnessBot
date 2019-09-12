package commands

import (
	"fmt"
	"github.com/Defman21/madnessBot/commands"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/Defman21/madnessBot/notifier"
	"github.com/franela/goreq"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct{}

func (c *Command) UseLua() bool {
	return false
}

func (c *Command) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	login := update.Message.CommandArguments()
	userName := update.Message.From.UserName

	req := goreq.Request{
		Uri: "https://api.twitch.tv/helix/users",
		QueryString: struct {
			Login string
		}{
			Login: login,
		},
	}
	oauth.AddHeadersUsing("twitch", &req)
	res, err := req.Do()

	if err != nil {
		common.Log.Error().Err(err).Msg("Request failed")
		return
	}

	type User struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	var user User

	err = res.Body.FromJsonTo(&user)
	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to unmarshal twitch response")
	}

	if len(user.Data) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такого пидора нет")
		_, _ = api.Send(msg)
		return
	}

	notifier.Get().Remove(user.Data[0].ID, fmt.Sprintf("@%s", userName))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("отписал тебя от оповещений %s", login))
	_, _ = api.Send(msg)
}

func init() {
	commands.Register("unnotify_me", &Command{})
}
