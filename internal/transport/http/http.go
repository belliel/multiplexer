package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/belliel/multiplexer/internal/transport/http/api"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultAddr               = ":8080"
	readTimeout               = 5 * time.Second
	writeTimeout              = 30 * time.Second
	idleConnectionTimeout     = 3 * time.Second
	defaultMaxConnectionLimit = 100
)

type Server struct {
	masterCtx         context.Context
	addr              string
	debug             bool
	connections       int32
	connectionsLimit  int32
	idleConnectionsCh chan struct{}
	instance          *http.Server
}

func NewServer(ctx context.Context, debug bool, addr string, connectionLimit int32) *Server {
	server := &Server{
		masterCtx:         ctx,
		addr:              addr,
		debug:             debug,
		connections:       0,
		connectionsLimit:  defaultMaxConnectionLimit,
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

	if server.connectionsLimit != connectionLimit {
		server.connectionsLimit = connectionLimit
	}

	server.getHandlers()

	return server
}

func (s *Server) throttleMiddleware(handler http.Handler) http.Handler {
	once := &sync.Once{}
	once.Do(func() {
		go func() {
			for {
				time.Sleep(1 * time.Second)
				fmt.Println(s.connections)
			}
		}()
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if s.connectionsLimit == 0 || atomic.LoadInt32(&s.connections) < s.connectionsLimit {
			atomic.AddInt32(&s.connections, 1)
			handler.ServeHTTP(w, r)
			atomic.AddInt32(&s.connections, -1)
			_ = ""
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("connections limit exceeded"))
		}
	})
}

func (s *Server) recoverMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			carry := recover()
			if carry != nil {
				var err error
				switch t := carry.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				log.Printf("[PANIC] [%s] %d %s\n", r.Method, http.StatusInternalServerError, err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) loggerMiddleware(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			status:         200,
		}

		start := time.Now()
		handler.ServeHTTP(recorder, r)
		end := time.Now().Sub(start)

		log.Printf(
			"[INFO] [%s] [%d] %s | %f sec | %s\n",
			r.Method, recorder.status, r.RemoteAddr, end.Seconds(), r.URL,
		)
	})
}

func (s *Server) getHandlers() {
	mux := http.DefaultServeMux

	mux.HandleFunc("/process/urls", api.ProcessUrls)

	middlewares := []func(http.Handler) http.Handler{
		s.throttleMiddleware,
		s.loggerMiddleware,
		s.recoverMiddleware,
	}

	s.instance.Handler = mux

	for _, middleware := range middlewares {
		s.instance.Handler = middleware(s.instance.Handler)
	}
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
	s.instance.SetKeepAlivesEnabled(false)
	close(s.idleConnectionsCh)
}

func (s *Server) WaitForGracefulShutdown() {
	<-s.idleConnectionsCh
}
