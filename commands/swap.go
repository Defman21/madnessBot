package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"madnessBot/common/helpers"
	"strings"
)

type SwapCmd struct{}

func (c SwapCmd) UseLua() bool {
	return false
}

func (c SwapCmd) Run(api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	dict := map[rune]rune{
		'q':  'й',
		'w':  'ц',
		'e':  'у',
		'r':  'к',
		't':  'е',
		'y':  'н',
		'u':  'г',
		'i':  'ш',
		'o':  'щ',
		'p':  'з',
		'[':  'х',
		']':  'ъ',
		'a':  'ф',
		's':  'ы',
		'd':  'в',
		'f':  'а',
		'g':  'п',
		'h':  'р',
		'j':  'о',
		'k':  'л',
		'l':  'д',
		';':  'ж',
		'\'': 'э',
		'z':  'я',
		'x':  'ч',
		'c':  'с',
		'v':  'м',
		'b':  'и',
		'n':  'т',
		'm':  'ь',
		',':  'б',
		'.':  'ю',
		'/':  '.',
	}
	text := update.Message.ReplyToMessage.Text
	fucked := []rune(strings.ToLower(text))
	for i, char := range fucked {
		val, ok := dict[char]
		if ok {
			fucked[i] = val
		} else {
			fucked[i] = char
		}
	}

	helpers.SendMessage(api, update, string(fucked), true, false)
}
