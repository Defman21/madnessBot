package notifier

import (
	"encoding/gob"
	"github.com/Defman21/madnessBot/common/logger"
	"os"
	"strings"
)

type gobStruct map[string][]string

type Notifier struct {
	userMap gobStruct
}

var singleton *Notifier

func Get() *Notifier {
	if singleton == nil {
		singleton = &Notifier{userMap: gobStruct{}}
		singleton.Load()
	}
	return singleton
}

const file = "./data/notifier.gob"

func (n *Notifier) Load() {
	if _, err := os.Stat(file); err == nil {
		file, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
		defer file.Close()

		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to load notifier file")
		}

		dec := gob.NewDecoder(file)
		err = dec.Decode(&n.userMap)

		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to decode notifier file")
		}

		logger.Log.Info().Interface("map", n.userMap).Msg("Loaded notifier map")
	} else if os.IsNotExist(err) {
		n.saveGob(&gobStruct{})
	}
}

func (n Notifier) saveGob(userMap *gobStruct) {
	file, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create notifier gob")
		return
	}

	enc := gob.NewEncoder(file)
	if userMap != nil {
		err = enc.Encode(userMap)
	} else {
		err = enc.Encode(n.userMap)
	}

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to save notifier gob")
	}
}

func (n *Notifier) Add(userID string, userName string) {
	if n.userMap[userID] == nil {
		n.userMap[userID] = []string{}
	}
	n.userMap[userID] = append(n.userMap[userID], userName)
	n.saveGob(nil)
}

func (n *Notifier) Remove(userID string, userName string) {
	var newList []string
	for _, user := range n.userMap[userID] {
		if user != userName {
			newList = append(newList, user)
		}
	}
	n.userMap[userID] = newList
	n.saveGob(nil)
}

func (n Notifier) GenerateNotifyString(userID string) string {
	return strings.Join(n.userMap[userID], ", ")
}
