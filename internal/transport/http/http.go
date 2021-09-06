package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	defaultAddr           = ":8080"
	readTimeout           = 5 * time.Second
	writeTimeout          = 30 * time.Second
	idleConnectionTimeout = 3 * time.Second
)

type Server struct {
	masterCtx         context.Context
	addr              string
	debug             bool
	idleConnectionsCh chan struct{}
	instance          *http.Server
}

func NewServer(ctx context.Context, debug bool, addr string) *Server {
	server := &Server{
		masterCtx: ctx,
		addr:      addr,
		debug:     debug,
		idleConnectionsCh: make(chan struct{}),
		instance: &http.Server{
			Addr:         defaultAddr,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}

	if server.instance.Addr != addr {
		server.instance.Addr = addr
	}

	return server
}

func (s *Server) Listen() error {
	go s.ListenCtxForGracefulShutdown()
	log.Printf("[INFO] Serving HTTP on \"%s\"", s.addr)
	return s.instance.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.instance.Shutdown(context.Background()) // masterContext cancels first
}

func (s *Server) ListenCtxForGracefulShutdown() {
	<-s.masterCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), idleConnectionTimeout)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] HTTP server Shutdown: %v", err)
	}

	log.Println("[INFO] Processing idle connections before termination")
	close(s.idleConnectionsCh)
}

func (s *Server) WaitForGracefulShutdown() {
	<-s.idleConnectionsCh
}
