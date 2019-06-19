package twitch

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/Defman21/madnessBot/common"
	"github.com/Defman21/madnessBot/common/oauth"
	"github.com/franela/goreq"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type twitchOauth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    time.Time
}

const file = "./data/twitch-state.gob"
const oauthUrl = "https://id.twitch.tv/oauth2/token"

var twitchInstance = &twitchOauth{}

func init() {
	oauth.Register("twitch", twitchInstance)
}

func (t *twitchOauth) Init() {
	if _, err := os.Stat(file); err == nil {
		file, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
		defer file.Close()

		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to load twitch state file")
		}

		dec := gob.NewDecoder(file)
		err = dec.Decode(t)

		if err != nil {
			common.Log.Error().Err(err).Msg("Failed to decode twitch state file")
		}

		common.Log.Info().Interface("state", twitchInstance).Msg("Loaded twitch oauth state")
	} else if os.IsNotExist(err) {
		t.Refresh()
	}
}

func (t *twitchOauth) Save() {
	file, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to open twitch state file")
		return
	}

	enc := gob.NewEncoder(file)

	if err = enc.Encode(t); err != nil {
		common.Log.Error().Err(err).Msg("Failed to save twitch state file")
		return
	}

	common.Log.Info().Interface("state", t).Msg("Saved twitch auth state")
}

func (t *twitchOauth) Authorize() {
	queryParams := url.Values{}
	queryParams.Add("client_id", os.Getenv("TWITCH_CLIENT_ID"))
	queryParams.Add("client_secret", os.Getenv("TWITCH_CLIENT_SECRET"))
	queryParams.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", oauthUrl, nil)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to create a request")
		return
	}

	req.URL.RawQuery = queryParams.Encode()

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	resp, err := client.Do(req)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to send the request")
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to read body of the request")
		return
	}

	err = json.Unmarshal(body, t)

	if err != nil {
		common.Log.Error().Err(err).Str("body", string(body)).Msg("Failed to parse JSON")
		return
	}

	t.UpdateExpire()

	common.Log.Info().Interface("oauth", t).Msg("Created tokens successfully")

	return
}

func (t *twitchOauth) Refresh() {
	t.Authorize()
	t.Save()
}

func (t *twitchOauth) UpdateExpire() {
	t.ExpiresAt = time.Now().Local().Add(time.Second * time.Duration(t.ExpiresIn))
}

func (t *twitchOauth) AddHeaders(request *goreq.Request) {
	request.AddHeader("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	request.AddHeader("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
}

func (t *twitchOauth) ExpiresSoon() bool {
	return time.Now().Local().After(t.ExpiresAt)
}