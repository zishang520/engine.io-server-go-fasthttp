package types

import (
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/zishang520/engine.io/v2/errors"
	"github.com/zishang520/engine.io/v2/events"
	"github.com/zishang520/engine.io/v2/utils"
)

type HttpServer struct {
	events.EventEmitter
	*ServeMux

	servers []any
	mu      sync.RWMutex
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

	s.servers = append(s.servers, server)

	return server
}

func (s *HttpServer) Close(fn func(error)) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.Emit("close")

	if s.servers != nil {
		var closingErr, serverErr error
		for _, server := range s.servers {
			switch srv := server.(type) {
			case *fasthttp.Server:
				serverErr = srv.Shutdown()
			default:
				serverErr = errors.New("unknown server type")
			}
			if serverErr != nil && closingErr == nil {
				closingErr = serverErr
			}
		}

		if closingErr != nil {
			err = fmt.Errorf("error occurred while closing servers: %v", closingErr)
		}
	}

	if fn != nil {
		defer fn(err)
	}

	return err
}

func (s *HttpServer) Listen(addr string, fn Callable) *fasthttp.Server {
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

func (s *HttpServer) ListenTLS(addr string, certFile string, keyFile string, fn Callable) *fasthttp.Server {
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
