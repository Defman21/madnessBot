package wiki

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"madnessBot/common/helpers"
	"net/url"
)

func HandleUpdate(api *tgbotapi.BotAPI, update *tgbotapi.Update, regexMatch []string) {
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
				Titles:        regexMatch[1],
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
		return
	}

	var text string
	if len(data.Query.Pages) != 0 && len(data.Query.Pages[0].Extract) != 0 {
		page := data.Query.Pages[0]
		text = fmt.Sprintf(
			"[%v](https://ru.wikipedia.org/wiki/%v)\n%v\n",
			helpers.EscapeMarkdownV2(page.Title),
			url.QueryEscape(page.Title),
			helpers.EscapeMarkdownV2(page.Extract),
		)
	} else {
		text = fmt.Sprintf("Вики не знает\\. [В гугл\\!](https://lmgtfy.com/?q=%s)", url.QueryEscape(regexMatch[1]))
	}
	helpers.SendMessage(api, update, text, true, true)
}
