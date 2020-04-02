package twitch

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"madnessBot/common/logger"
	"madnessBot/config"
	"madnessBot/redis"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const oauthUrl = "https://id.twitch.tv/oauth2/token"
const redisKey = "madnessBot:state:oauth:twitch"

type twitchOauth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    time.Time
}

var Instance = &twitchOauth{}

func (t twitchOauth) getRedisMap() map[string]interface{} {
	return map[string]interface{}{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
		"expires_in":    strconv.FormatInt(t.ExpiresIn, 10),
		"expires_at":    t.ExpiresAt.Format(time.RFC3339),
	}
}

func (t *twitchOauth) setFromRedisMap(redisMap map[string]string) {
	t.AccessToken = redisMap["access_token"]
	t.RefreshToken = redisMap["refresh_token"]
	i, _ := strconv.ParseInt(redisMap["expires_in"], 10, 64)
	t.ExpiresIn = i
	t.ExpiresAt, _ = time.Parse(time.RFC3339, redisMap["expires_at"])
}

func (t *twitchOauth) Init() {
	_redis := redis.Get()
	existsInt, err := _redis.Exists(redisKey).Result()
	if err != nil {
		logger.Log.Error().Err(err).Str("key", redisKey).Msg("Failed to EXISTS redis key")
	}

	if existsInt == 0 {
		t.Refresh()
		return
	}

	redisMap, err := redis.Get().HGetAll(redisKey).Result()
	if err != nil {
		logger.Log.Error().Err(err).Str("key", redisKey).Msg("Failed to HGETALL redis key")
	}

	t.setFromRedisMap(redisMap)
	logger.Log.Info().Interface("state", Instance).Msg("Loaded twitch oauth state")
}

func (t *twitchOauth) Save() {
	fields := t.getRedisMap()
	_, err := redis.Get().HSet(redisKey, fields).Result()
	if err != nil {
		logger.Log.Error().Err(err).
			Str("key", redisKey).
			Fields(fields).
			Msg("Failed to HSET redis key")
	}
	logger.Log.Info().Interface("state", t).Msg("Saved twitch auth state")
}

func (t *twitchOauth) Authorize() {
	queryParams := url.Values{}
	queryParams.Add("client_id", config.Config.Twitch.ClientID)
	queryParams.Add("client_secret", config.Config.Twitch.ClientSecret)
	queryParams.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", oauthUrl, nil)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create a request")
		return
	}

	req.URL.RawQuery = queryParams.Encode()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to send the request")
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to read body of the request")
		return
	}

	err = json.Unmarshal(body, t)

	if err != nil {
		logger.Log.Error().Err(err).Str("body", string(body)).Msg("Failed to parse JSON")
		return
	}

	t.UpdateExpire()

	logger.Log.Info().Interface("oauth", t).Msg("Created tokens successfully")

	return
}

func (t *twitchOauth) Refresh() {
	t.Authorize()
	t.Save()
}

func (t *twitchOauth) UpdateExpire() {
	t.ExpiresAt = time.Now().Local().Add(time.Second * time.Duration(t.ExpiresIn))
}

func (t *twitchOauth) AddHeaders(agent *gorequest.SuperAgent) {
	agent.Set("Client-UserID", config.Config.Twitch.ClientID)
	agent.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
}

func (t *twitchOauth) ExpiresSoon() bool {
	return time.Now().Local().After(t.ExpiresAt.Add(-1 * 12 * time.Hour))
}
