package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"madnessBot/common/logger"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

var log = &logger.Log

type twitchConfig struct {
	ClientID     string        `yaml:"client_id"`
	ClientSecret string        `yaml:"client_secret"`
	Webhook      webhookConfig `yaml:"webhook"`
	Enabled      bool          `yaml:"enabled"`
}

type graphiteConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

type newsConfig struct {
	Enabled bool             `yaml:"enabled"`
	Token   string           `yaml:"token"`
	Sources map[string]int64 `yaml:"sources"`
}

type webhookConfig struct {
	Enable bool   `yaml:"enabled"`
	URL    string `yaml:"url"`
	Path   string `yaml:"path"`
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
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
}

func (cfg serverConfig) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

type redisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type config struct {
	Token            string          `yaml:"token"`
	Twitch           *twitchConfig   `yaml:"twitch"`
	ChatID           int64           `yaml:"chat_id"`
	BoostyChatID     int64           `yaml:"boosty_chat_id"`
	ErrorChatID      int64           `yaml:"error_chat_id"`
	BoostyLink       string          `yaml:"boosty_link"`
	LogLevel         string          `yaml:"log_level"`
	Graphite         *graphiteConfig `yaml:"graphite"`
	NotificationsLRU int             `yaml:"notifications_lru_cache"`
	Admins           []int64         `yaml:"admins"`
	News             *newsConfig     `yaml:"news"`
	Webhook          *webhookConfig  `yaml:"webhook"`
	Server           serverConfig    `yaml:"server"`
	Redis            *redisConfig    `yaml:"redis"`
	MessageThreshold int64           `yaml:"message_threshold"`
}

func (c config) GetAdmins() map[int64]bool {
	ret := map[int64]bool{}
	for _, id := range c.Admins {
		ret[id] = true
	}
	return ret
}

const configName = "config.yml"

func getCwd() string {
	ex, err := os.Executable()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current executable")
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
	log.Info().Str("cfg", getConfigPath()).Msg("Config")
	file, err := os.OpenFile(getConfigPath(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open config.yml")
		return false
	}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode config.yml")
		return false
	}
	Initialized <- true
	return true
}
