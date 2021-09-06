package http

import (
	"context"
	"net/http"
)

type Server struct {
	masterCtx context.Context
	port      string
	debug     bool
	instance  *http.Server
}

func NewServer(ctx context.Context, debug bool, port string) *Server {
	server := &Server{
		masterCtx: ctx,
		port:      port,
		debug:     debug,
	}

	server.instance = &http.Server{
		Addr:              ":" + port,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	return server
}

func (s *Server) Listen() {
	panic("implement me")
}

func (s *Server) Shutdown() {
	panic("implement me")
}
