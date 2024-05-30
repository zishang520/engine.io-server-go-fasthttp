package types

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/zishang520/engine.io/v2/errors"
	"github.com/zishang520/engine.io/v2/events"
	_types "github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
)

type HttpServer struct {
	events.EventEmitter
	*ServeMux

	servers *_types.Slice[any]
}

func NewWebServer(defaultHandler Handler) *HttpServer {

	s := &HttpServer{
		EventEmitter: events.New(),
		ServeMux:     NewServeMux(defaultHandler),
	}
	return s
}

// Deprecated: this method will be removed in the next major release, please use [NewWebServer] instead.
func CreateServer(defaultHandler Handler) *HttpServer {
	return NewWebServer(defaultHandler)
}

func (s *HttpServer) httpServer(handler Handler) *fasthttp.Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	server := &fasthttp.Server{
		Handler: handler.FastHTTP,
		Logger:  utils.Log(),
		ErrorHandler: func(ctx *fasthttp.RequestCtx, err error) {
			s.Emit("error", err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		}}

	s.servers.Push(server)

	return server
}

func (s *HttpServer) Close(fn func(error)) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.Emit("close")

	if s.servers != nil {
		var closingErr, serverErr error
		s.servers.Range(func(server any, _ int) bool {
			switch srv := server.(type) {
			case *fasthttp.Server:
				serverErr = srv.Shutdown()
			default:
				serverErr = errors.New("unknown server type")
			}
			if serverErr != nil && closingErr == nil {
				closingErr = serverErr
			}
			return true
		})

		if closingErr != nil {
			err = fmt.Errorf("error occurred while closing servers: %v", closingErr)
		}
	}

	if fn != nil {
		defer fn(err)
	}

	return err
}

func (s *HttpServer) Listen(addr string, fn _types.Callable) *fasthttp.Server {
	server := s.httpServer(s)
	go func() {
		if err := server.ListenAndServe(addr); err != nil {
			panic(err)
		}
	}()

	if fn != nil {
		defer fn()
	}
	s.Emit("listening")

	return server
}

func (s *HttpServer) ListenTLS(addr string, certFile string, keyFile string, fn _types.Callable) *fasthttp.Server {
	server := s.httpServer(s)
	go func() {
		if err := server.ListenAndServeTLS(addr, certFile, keyFile); err != nil {
			panic(err)
		}
	}()

	if fn != nil {
		defer fn()
	}
	s.Emit("listening")

	return server
}
