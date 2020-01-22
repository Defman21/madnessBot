package online_state

import (
	"encoding/gob"
	"github.com/Defman21/madnessBot/common/logger"
	"os"
)

var state = map[string]bool{}

func Push(username string, isOnline bool) {
	state[username] = isOnline
	saveState()
}

func GetOnline() (result []string) {
	for username, isOnline := range state {
		if isOnline {
			result = append(result, username)
		}
	}
	return
}
func saveState() {
	if file, err := os.OpenFile(
		"./data/online-state.gob",
		os.O_CREATE|os.O_WRONLY, os.ModePerm,
	); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to open/create online-state.gob")
	} else {
		defer file.Close()

		encoder := gob.NewEncoder(file)
		err = encoder.Encode(state)

		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to encode data")
		}

		logger.Log.Info().Interface("state", state).Msg("Saved online state")
	}
}

func loadState() {
	if file, err := os.OpenFile("./data/online-state.gob", os.O_RDONLY, os.ModePerm); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to open online-state.gob")
	} else {
		defer file.Close()

		decoder := gob.NewDecoder(file)

		err = decoder.Decode(&state)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to decode online-state.gob")
		}

		logger.Log.Info().Interface("state", state).Msg("Loaded online state")
	}
}

func init() {
	loadState()
}
