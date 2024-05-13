package types

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"github.com/zishang520/engine.io/v2/errors"
	"github.com/zishang520/engine.io/v2/events"
	"github.com/zishang520/engine.io/v2/utils"
)

type HttpContext struct {
	events.EventEmitter

	Websocket *WebSocketConn

	Cleanup Callable

	requestCtx *fasthttp.RequestCtx

	headers *utils.ParameterBag
	query   *utils.ParameterBag

	method      string
	pathInfo    string
	isHostValid bool

	isDone bool
	done   chan Void
	mu     sync.RWMutex

	statusCode      int
	mu_wh           sync.RWMutex
	ResponseHeaders *utils.ParameterBag

	mu_w sync.Mutex
}

func NewHttpContext(ctx *fasthttp.RequestCtx) *HttpContext {
	c := &HttpContext{}
	c.EventEmitter = events.New()
	c.done = make(chan Void)

	c.requestCtx = ctx

	c.headers = utils.NewParameterBag(nil)
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		c.headers.Set(strconv.B2S(key), strconv.B2S(value))
	})

	c.query = utils.NewParameterBag(nil)
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		c.query.Set(strconv.B2S(key), strconv.B2S(value))
	})

	c.isHostValid = true

	c.ResponseHeaders = utils.NewParameterBag(nil)
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isDone {
		close(c.done)
		c.isDone = true
	}
}

func (c *HttpContext) Done() <-chan Void {
	return c.done
}

func (c *HttpContext) IsDone() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.isDone
}

func (c *HttpContext) SetStatusCode(statusCode int) {
	c.mu_wh.Lock()
	defer c.mu_wh.Unlock()

	c.statusCode = statusCode
}

func (c *HttpContext) GetStatusCode() int {
	c.mu_wh.RLock()
	defer c.mu_wh.RUnlock()

	return c.statusCode
}

func (c *HttpContext) Write(wb []byte) (int, error) {
	c.mu_w.Lock()
	defer c.mu_w.Unlock()

	if !c.IsDone() {
		defer c.Flush()

		for k, v := range c.ResponseHeaders.All() {
			c.requestCtx.Response.Header.Set(k, v[0])
		}
		c.requestCtx.SetStatusCode(c.GetStatusCode())

		return c.requestCtx.Write(wb)
	}
	return 0, errors.New("You cannot write data repeatedly.").Err()
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
	if c.method != "" {
		return c.method
	}

	c.method = strings.ToUpper(strconv.B2S(c.requestCtx.Method()))
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
	// trim and remove port number from host
	// host is lowercase as per RFC 952/2181
	host = regexp.MustCompile(`:\d+$`).ReplaceAllString(strings.TrimSpace(host), "")
	// as the host can come from the user (HTTP_HOST and depending on the configuration, SERVER_NAME too can come from the user)
	// check that it does not contain forbidden characters (see RFC 952 and RFC 2181)
	if host != "" {
		if host = regexp.MustCompile(`(?:^\[)?[a-zA-Z0-9-:\]_]+\.?`).ReplaceAllString(host, ""); host != "" {
			if !c.isHostValid {
				return "", nil
			}
			c.isHostValid = false
			return "", errors.New(fmt.Sprintf(`Invalid Host "%s".`, host)).Err()
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
