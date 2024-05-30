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

var _transports map[string]*transports = map[string]*transports{
	"polling": {
		// Polling polymorphic New.
		New: func(ctx *types.HttpContext) Transport {
			if ctx.Query().Has("j") {
				return NewJSONP(ctx)
			}
			return NewPolling(ctx)
		},
		HandlesUpgrades: false,
		UpgradesTo:      _types.NewSet("websocket"),
	},

	"websocket": {
		New: func(ctx *types.HttpContext) Transport {
			return NewWebSocket(ctx)
		},
		HandlesUpgrades: true,
		UpgradesTo:      _types.NewSet[string](),
	},
}

func Transports() map[string]*transports {
	return _transports
}
