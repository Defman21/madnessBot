package common

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/franela/goreq"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TwitchOauth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    time.Time
}

const File = "./data/twitch-state.gob"
const OauthUrl = "https://id.twitch.tv/oauth2/token"

func (o *TwitchOauth) Load() error {
	if _, err := os.Stat(File); err == nil {
		file, err := os.OpenFile(File, os.O_RDONLY, os.ModePerm)
		defer file.Close()

		if err != nil {
			Log.Error().Err(err).Msg("Failed to load twitch state file")
			return err
		}

		dec := gob.NewDecoder(file)
		err = dec.Decode(o)

		if err != nil {
			Log.Error().Err(err).Msg("Failed to decode twitch state file")
			return err
		}

		Log.Info().Interface("state", o).Msg("Loaded twitch oauth state")
	} else if os.IsNotExist(err) {
		if err = o.Authorize(); err != nil {
			Log.Error().Err(err).Msg("Failed to authorize")
			return err
		}

		if err = o.Save(); err != nil {
			Log.Error().Err(err).Msg("Failed to save twitch state")
		}

		return nil
	} else {
		return err
	}

	return nil
}

func (o *TwitchOauth) Save() error {
	file, err := os.OpenFile(File, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	if err != nil {
		Log.Error().Err(err).Msg("Failed to open twitch state file")
		return err
	}

	enc := gob.NewEncoder(file)

	if err = enc.Encode(o); err != nil {
		Log.Error().Err(err).Msg("Failed to save twitch state file")
		return err
	}

	Log.Info().Interface("state", o).Msg("Saved twitch auth state")
	return nil
}

func (o *TwitchOauth) Authorize() error {
	queryParams := url.Values{}
	queryParams.Add("client_id", os.Getenv("TWITCH_CLIENT_ID"))
	queryParams.Add("client_secret", os.Getenv("TWITCH_CLIENT_SECRET"))
	queryParams.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", OauthUrl, nil)

	if err != nil {
		Log.Error().Err(err).Msg("Failed to create a request")
		return err
	}

	req.URL.RawQuery = queryParams.Encode()

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	resp, err := client.Do(req)

	if err != nil {
		Log.Error().Err(err).Msg("Failed to send the request")
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		Log.Error().Err(err).Msg("Failed to read body of the request")
		return err
	}

	err = json.Unmarshal(body, o)

	if err != nil {
		Log.Error().Err(err).Str("body", string(body)).Msg("Failed to parse JSON")
		return err
	}

	o.UpdateExpire()

	Log.Info().Interface("oauth", o).Msg("Created tokens successfully")

	return nil
}

func (o *TwitchOauth) Refresh() {
	_ = o.Authorize()
	_ = o.Save()
}

func (o *TwitchOauth) UpdateExpire() {
	o.ExpiresAt = time.Now().Local().Add(time.Second * time.Duration(o.ExpiresIn))
}

func (o *TwitchOauth) AddHeaders(request *goreq.Request) {
	request.AddHeader("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	request.AddHeader("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))
}

var TwitchOauthState *TwitchOauth

func init() {
	TwitchOauthState = &TwitchOauth{}
}
