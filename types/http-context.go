package types

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"github.com/zishang520/engine.io/v2/errors"
	"github.com/zishang520/engine.io/v2/events"
	_types "github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
)

type HttpContext struct {
	events.EventEmitter

	Websocket *WebSocketConn

	Cleanup _types.Callable

	requestCtx *fasthttp.RequestCtx

	headers *utils.ParameterBag
	query   *utils.ParameterBag

	method      string
	pathInfo    string
	isHostValid bool

	isDone atomic.Bool
	done   chan _types.Void

	statusCode      atomic.Value
	ResponseHeaders *utils.ParameterBag

	mu sync.Mutex
}

func NewHttpContext(ctx *fasthttp.RequestCtx) *HttpContext {
	c := &HttpContext{
		EventEmitter:    events.New(),
		done:            make(chan _types.Void),
		requestCtx:      ctx,
		headers:         utils.NewParameterBag(nil),
		query:           utils.NewParameterBag(nil),
		isHostValid:     true,
		ResponseHeaders: utils.NewParameterBag(nil),
	}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		c.headers.Set(strconv.B2S(key), strconv.B2S(value))
	})
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		c.query.Set(strconv.B2S(key), strconv.B2S(value))
	})
	ctx.Response.Header.VisitAll(func(key, value []byte) {
		c.ResponseHeaders.Set(strconv.B2S(key), strconv.B2S(value))
	})

	go func() {
		select {
		case <-c.done:
			c.Emit("close")
		}
	}()

	return c
}

func (c *HttpContext) Flush() {
	if c.isDone.CompareAndSwap(false, true) {
		close(c.done)
	}
}

func (c *HttpContext) Done() <-chan _types.Void {
	return c.done
}

func (c *HttpContext) IsDone() bool {
	return c.isDone.Load()
}

func (c *HttpContext) SetStatusCode(statusCode int) {
	c.statusCode.Store(statusCode)
}

func (c *HttpContext) GetStatusCode() int {
	if v, ok := c.statusCode.Load().(int); ok {
		return v
	}
	return fasthttp.StatusOK
}

func (c *HttpContext) Write(wb []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.IsDone() {
		return 0, errors.New("you cannot write data repeatedly").Err()
	}
	defer c.Flush()

	for k, v := range c.ResponseHeaders.All() {
		c.requestCtx.Response.Header.Set(k, v[0])
	}
	c.requestCtx.SetStatusCode(c.GetStatusCode())

	return c.requestCtx.Write(wb)
}

func (c *HttpContext) RequestCtx() *fasthttp.RequestCtx {
	return c.requestCtx
}

func (c *HttpContext) Headers() *utils.ParameterBag {
	return c.headers
}

func (c *HttpContext) Query() *utils.ParameterBag {
	return c.query
}

func (c *HttpContext) GetPathInfo() string {
	if c.pathInfo == "" {
		c.pathInfo = strconv.B2S(c.requestCtx.Path())
	}
	return c.pathInfo
}

func (c *HttpContext) Get(key string, _default ...string) string {
	v, _ := c.query.Get(key, _default...)
	return v
}

func (c *HttpContext) Gets(key string, _default ...[]string) []string {
	v, _ := c.query.Gets(key, _default...)
	return v
}

func (c *HttpContext) Method() string {
	return c.GetMethod()
}

func (c *HttpContext) GetMethod() string {
	if c.method == "" {
		c.method = strings.ToUpper(strconv.B2S(c.requestCtx.Method()))
	}
	return c.method
}

func (c *HttpContext) Path() string {
	if pattern := strings.Trim(c.GetPathInfo(), "/"); pattern != "" {
		return pattern
	}
	return "/"
}

func (c *HttpContext) GetHost() (string, error) {
	host := strconv.B2S(c.requestCtx.Host())
	host = regexp.MustCompile(`:\d+$`).ReplaceAllString(host, "")

	if host != "" {
		if host = regexp.MustCompile(`(?:^\[)?[a-zA-Z0-9-:\]_]+\.?`).ReplaceAllString(host, ""); host != "" {
			if !c.isHostValid {
				return "", nil
			}
			c.isHostValid = false
			return "", errors.New(fmt.Sprintf(`Invalid host "%s".`, host)).Err()
		}
	}
	return host, nil
}

func (c *HttpContext) UserAgent() string {
	return c.headers.Peek("User-Agent")
}

func (c *HttpContext) Secure() bool {
	return c.requestCtx.IsTLS()
}
