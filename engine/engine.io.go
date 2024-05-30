package engine

import (
	"github.com/valyala/fasthttp"
	"github.com/zishang520/engine.io-server-go-fasthttp/v2/types"
	_types "github.com/zishang520/engine.io/v2/types"
)

const Protocol = 4

func New(server any, args ...any) Server {
	switch s := server.(type) {
	case *types.HttpServer:
		return Attach(s, append(args, nil)[0])
	case any:
		return NewServer(s)
	}
	return NewServer(nil)
}

// Creates an fasthttp.Server exclusively used for WS upgrades.
func Listen(addr string, options any, fn _types.Callable) Server {
	server := types.NewWebServer(types.HandlerFunc(func(ctx *fasthttp.RequestCtx) {
		ctx.Error("Not Implemented", fasthttp.StatusNotImplemented)
	}))

	// create engine server
	engine := Attach(server, options)
	engine.SetHttpServer(server)

	server.Listen(addr, fn)

	return engine
}

// Captures upgrade requests for a types.HttpServer.
func Attach(server *types.HttpServer, options any) Server {
	engine := NewServer(options)
	engine.Attach(server, options)
	return engine
}
