package transports

import (
	"github.com/zishang520/engine.io-server-go-fasthttp/v2/types"
	_types "github.com/zishang520/engine.io/v2/types"
)

type transports struct {
	New             func(*types.HttpContext) Transport
	HandlesUpgrades bool
	UpgradesTo      *_types.Set[string]
}

var _transports map[string]*transports

func init() {
	_transports = map[string]*transports{
		"websocket": {
			New: func(ctx *types.HttpContext) Transport {
				return NewWebSocket(ctx)
			},
			HandlesUpgrades: true,
			UpgradesTo:      _types.NewSet[string](),
		},
	}
}

func Transports() map[string]*transports {
	return _transports
}
