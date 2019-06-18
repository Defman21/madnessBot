package common

import (
	"encoding/gob"
	"os"
	"time"
)

type ResubscribeStateSingleton struct {
	ExpiresAt time.Time
}

func (r *ResubscribeStateSingleton) Save() {
	if file, err := os.OpenFile(
		"./data/resub-state.gob",
		os.O_CREATE|os.O_WRONLY, os.ModePerm,
	); err != nil {
		Log.Error().Err(err).Msg("Failed to open/create resub-state.gob")
	} else {
		defer file.Close()

		r.ExpiresAt = time.Now().Local().Add(time.Hour * time.Duration(8*24))

		encoder := gob.NewEncoder(file)
		err = encoder.Encode(r)

		if err != nil {
			Log.Error().Err(err).Msg("Failed to encode data")
		}

		Log.Info().Interface("state", r).Msg("Saved resubscribe state")
	}
}

func (r *ResubscribeStateSingleton) Load() {
	if file, err := os.OpenFile("./data/resub-state.gob", os.O_RDONLY, os.ModePerm); err != nil {
		Log.Error().Err(err).Msg("Failed to open resub-state.gob")
	} else {
		defer file.Close()

		decoder := gob.NewDecoder(file)

		err = decoder.Decode(r)
		if err != nil {
			Log.Error().Err(err).Msg("Failed to decode resub-state.gob")
		}

		Log.Info().Interface("state", r).Msg("Loaded resubscribe state")
	}
}

var ResubscribeState *ResubscribeStateSingleton

func init() {
	ResubscribeState = &ResubscribeStateSingleton{}
}
