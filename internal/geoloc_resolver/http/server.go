package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

type Server struct {
	server *http.Server

	done chan struct{}
}

// NewServer creates new HTTP server, that listens on addr and handles resolve
// IP requests using handler h.
func NewServer(h *Handler, addr string) *Server {
	if h == nil {
		log.Fatalln("[F] http: got nil handler for new server")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/resolve_ip", h.handleResolveIPRequest)

	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		done: make(chan struct{}),
	}
}

// Init initializes HTTP server.
func (s *Server) Init() {
	go func() {
		defer close(s.done)

		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[F] http: error on listen and serve: %s", err)
		}
	}()
}

// Shutdown gracefully shutdowns HTTP server.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	if err := s.server.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("[F] http: failed to shutdown: %s", err)
	}

	cancel()
	<-s.done
}
