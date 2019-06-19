package oauth

import (
	"github.com/Defman21/madnessBot/common"
	"github.com/franela/goreq"
)

type Interface interface {
	Init()
	Refresh()
	ExpiresSoon() bool
	AddHeaders(request *goreq.Request)
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

func AddHeadersUsing(name string, request *goreq.Request) {
	if instance := Get(name); instance != nil {
		instance.AddHeaders(request)
	}
}

func RefreshExpired() {
	for name, provider := range providers {
		if provider.ExpiresSoon() {
			provider.Refresh()
			common.Log.Info().Str("name", name).Msg("Refreshed provider")
		}
	}
}
