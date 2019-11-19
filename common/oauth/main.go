package oauth

import (
	"github.com/Defman21/madnessBot/common/logger"
	"github.com/parnurzeal/gorequest"
)

type Interface interface {
	Init()
	Refresh()
	ExpiresSoon() bool
	AddHeaders(agent *gorequest.SuperAgent)
}

var providers = make(map[string]Interface)

func Register(name string, instance Interface) {
	providers[name] = instance
	go instance.Init()
}

func Get(name string) Interface {
	if instance, ok := providers[name]; ok {
		return instance
	} else {
		return nil
	}
}

func AddHeadersUsing(name string, agent *gorequest.SuperAgent) {
	if instance := Get(name); instance != nil {
		instance.AddHeaders(agent)
	}
}

func RefreshExpired() {
	for name, provider := range providers {
		if provider.ExpiresSoon() {
			provider.Refresh()
			logger.Log.Info().Str("name", name).Msg("Refreshed provider")
		}
	}
}
