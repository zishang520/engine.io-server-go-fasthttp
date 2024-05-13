package types

import (
	"github.com/fasthttp/websocket"
	"github.com/zishang520/engine.io/v2/events"
)

type WebSocketConn struct {
	events.EventEmitter
	*websocket.Conn

	exit chan Void
}

func (t *WebSocketConn) Close() error {
	defer t.Emit("close")
	return t.Conn.Close()
}

func MakeWebSocketConn() *WebSocketConn {
	c := &WebSocketConn{
		EventEmitter: events.New(),

		exit: make(chan Void),
	}

	return c
}

func NewWebSocketConn() *WebSocketConn {
	t := MakeWebSocketConn()

	t.Construct()

	return t
}

// WebSocketConn Construct.
func (t *WebSocketConn) Construct() {
	t.Once("close", func(...any) {
		close(t.exit)
	})
}

func (t *WebSocketConn) Done() <-chan Void {
	return t.exit
}
