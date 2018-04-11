package commands

import (
	"github.com/Defman21/madnessBot/common"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
)

// Swap luuul
func Swap(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
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
	defer func() {
		if err := recover(); err != nil {
			common.Log.Warn(err)
		}
	}()
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

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(fucked))
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}
