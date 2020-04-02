package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Defman21/madnessBot/common/logger"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type twitchConfig struct {
	ClientID     string        `toml:"client_id"`
	ClientSecret string        `toml:"client_secret"`
	Webhook      webhookConfig `toml:"webhook"`
	Enabled      bool          `toml:"enabled"`
}

type graphiteConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Enabled bool   `toml:"enabled"`
}

type newsConfig struct {
	Enabled bool             `toml:"enabled"`
	Token   string           `toml:"token"`
	Sources map[string]int64 `toml:"sources"`
}

type webhookConfig struct {
	Enable bool   `toml:"enabled"`
	URL    string `toml:"url"`
	Path   string `toml:"path"`
}

func (cfg webhookConfig) GetURL(paths ...string) string {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Invalid URL")
		return ""
	}
	join := append([]string{u.Path, cfg.Path}, paths...)
	u.Path = path.Join(join...)
	return u.String()
}

func (cfg *webhookConfig) Enabled() bool {
	return cfg != nil && cfg.Enable
}

type serverConfig struct {
	Host string `toml:"host"`
	Port int64  `toml:"port"`
}

func (cfg serverConfig) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

type config struct {
	Token            string          `toml:"token"`
	Twitch           *twitchConfig   `toml:"twitch"`
	ChatID           int64           `toml:"chat_id"`
	BoostyChatID     int64           `toml:"boosty_chat_id"`
	ErrorChatID      int64           `toml:"error_chat_id"`
	CardNumber       string          `toml:"card_number"`
	LogLevel         string          `toml:"log_level"`
	Graphite         *graphiteConfig `toml:"graphite"`
	NotificationsLRU int             `toml:"notifications_lru_cache"`
	Admins           map[string]bool `toml:"admins"`
	Payers           map[string]bool `toml:"payers"`
	News             *newsConfig     `toml:"news"`
	Webhook          *webhookConfig  `toml:"webhook"`
	Server           serverConfig    `toml:"server"`
}

func (c config) GetAdmins() map[int64]bool {
	ret := map[int64]bool{}
	for k, v := range c.Admins {
		i, _ := strconv.ParseInt(k, 10, 64)
		ret[i] = v
	}
	return ret
}

func (c config) GetPayers() map[int64]bool {
	ret := map[int64]bool{}
	for k, v := range c.Payers {
		i, _ := strconv.ParseInt(k, 10, 64)
		ret[i] = v
	}
	return ret
}

func (c *config) AddPayer(payerID int64) {
	c.Payers[strconv.FormatInt(payerID, 10)] = true
}

func (c *config) RemovePayer(payerID int64) {
	delete(c.Payers, strconv.FormatInt(payerID, 10))
}

const configName = "config.toml"

func getCwd() string {
	ex, err := os.Executable()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get current executable")
		return ""
	}
	return filepath.Dir(ex)
}

func getConfigPath() string {
	return filepath.Join(getCwd(), configName)
}

var Config config
var Initialized = make(chan bool, 1)

func Init() bool {
	logger.Log.Info().Str("cfg", getConfigPath()).Msg("Config")
	if _, err := toml.DecodeFile(getConfigPath(), &Config); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to decode config.toml")
		return false
	}
	Initialized <- true
	return true
}

func Save() {
	logger.Log.Info().Str("cfg", getConfigPath()).Msg("Saving config")
	file, err := os.OpenFile(getConfigPath(), os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to open config file")
	}
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(Config)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to encode config")
	}
	err = file.Close()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to close config file")
	}
}
