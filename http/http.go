// Package http implements the HTTP transmission layer.
package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/christgf/ports"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/errgroup"
)

// PortService provides the business logic implementation for our HTTP API.
type PortService interface {
	StorePort(ctx context.Context, p ports.Port) error
	GetPortByID(ctx context.Context, portID string) (*ports.Port, error)
}

// Server is our HTTP server, a thin wrapper over the standard library.
type Server struct {
	server *http.Server
	logger *log.Logger

	Ports PortService
}

// NewServer creates and returns a new Server backed by the PortService provided.
func NewServer(ps PortService) *Server {
	mux := httprouter.New()
	srv := &Server{
		server: &http.Server{
			Handler: mux,
		},
		logger: log.New(os.Stderr, "http ", log.LstdFlags),
		Ports:  ps,
	}

	// HTTP API routes.
	{
		// Health and readiness probes.
		mux.HandlerFunc(http.MethodGet, "/alive", srv.HandleAlive)
		mux.HandlerFunc(http.MethodGet, "/ready", srv.HandleReady)

		// Ports API.
		mux.HandlerFunc(http.MethodGet, "/ports", srv.HandleGetPort)
		mux.HandlerFunc(http.MethodPost, "/ports", srv.HandleStorePort)
	}

	return srv
}

// Serve begins listening and serving incoming HTTP requests. It will block
// serving requests until the context is canceled, at which point it attempts to
// shut down gracefully without interrupting any active connections.
func (s *Server) Serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-ctx.Done()
		return s.server.Shutdown(context.Background())
	})
	g.Go(func() error {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	return g.Wait()
}

// Reply to an HTTP request with the specified HTTP code and an optional payload.
// The function does not otherwise end the request, the caller should ensure no
// further writes are done to w.
func (s *Server) Reply(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			s.logger.Printf("json.Encode: %v", err)
		}
	}
}

// ErrorResponse is the HTTP response body delivered for every HTTP request that
// cannot be processed.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ReplyErr examines the error provided and replies to an HTTP request with an
// appropriate HTTP code and payload. The function does not otherwise end the
// request; the caller should ensure no further writes are done to w.
func (s *Server) ReplyErr(w http.ResponseWriter, err error) {
	const (
		code = 500
	)
	s.logger.Printf("HTTP %d (%s): %v", code, http.StatusText(code), err)
	s.Reply(w, code, &ErrorResponse{
		Code:    strings.ToLower(http.StatusText(code)),
		Message: err.Error(),
	})
}
