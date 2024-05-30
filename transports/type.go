package transports

import (
	"github.com/zishang520/engine.io-go-parser/packet"
	"github.com/zishang520/engine.io-go-parser/parser"
	_types "github.com/zishang520/engine.io-go-parser/types"
	"github.com/zishang520/engine.io-server-go-fasthttp/v2/types"
	"github.com/zishang520/engine.io/v2/events"
	e_types "github.com/zishang520/engine.io/v2/types"
)

type (
	Transport interface {
		// #extends

		events.EventEmitter

		// #prototype

		Prototype(Transport)
		Proto() Transport

		// #setters

		SetSid(string)
		SetWritable(bool)
		SetReq(*types.HttpContext)
		SetSupportsBinary(bool)
		SetReadyState(string)
		SetHttpCompression(*e_types.HttpCompression)
		SetPerMessageDeflate(*e_types.PerMessageDeflate)
		SetMaxHttpBufferSize(int64)

		// #getters

		Sid() string
		Writable() bool
		Protocol() int
		// @protected
		Discarded() bool
		// @protected
		Parser() parser.Parser
		// @protected
		Req() *types.HttpContext
		// @protected
		SupportsBinary() bool
		ReadyState() string
		HttpCompression() *e_types.HttpCompression
		PerMessageDeflate() *e_types.PerMessageDeflate
		MaxHttpBufferSize() int64
		// @abstract
		HandlesUpgrades() bool
		// @abstract
		SupportsFraming() bool
		// @abstract
		Name() string

		// #methods

		// Construct() should be called after calling Prototype()
		Construct(*types.HttpContext)
		// @private
		// Flags the transport as discarded.
		Discard()
		// @protected
		// Called with an incoming HTTP request.
		OnRequest(*types.HttpContext)
		// @private
		// Closes the transport.
		Close(...e_types.Callable)
		// @protected
		// Called with a transport error.
		OnError(string, error)
		// @protected
		// Called with parsed out a packets from the data stream.
		OnPacket(*packet.Packet)
		// @protected
		// Called with the encoded packet data.
		OnData(_types.BufferInterface)
		// @protected
		// Called upon transport close.
		OnClose()
		// @protected
		// @abstract
		// Writes a packet payload.
		Send([]*packet.Packet)
		// @protected
		// @abstract
		// Closes the transport.
		DoClose(e_types.Callable)
	}

	Polling interface {
		// #extends

		Transport

		// #methods

		DoWrite(*types.HttpContext, _types.BufferInterface, *packet.Options, func(*types.HttpContext))
	}

	Jsonp interface {
		// #extends

		Polling
	}

	Websocket interface {
		// #extends

		Transport
	}
)
