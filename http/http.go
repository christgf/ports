// Package http implements the HTTP transmission layer.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/christgf/ports"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/errgroup"
)

// PortService provides the business logic implementation for our HTTP API.
type PortService interface {
	StorePort(ctx context.Context, p ports.Port) error
	GetPortByID(ctx context.Context, portID string) (*ports.Port, error)
}

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 60 * time.Second
	defaultIdleTimeout  = 240 * time.Second
)

// Server is our HTTP server, a thin wrapper over the standard library.
type Server struct {
	server *http.Server
	logger *log.Logger

	Addr  string
	Ports PortService
}

// WithReadTimeout specifies how long the server should wait for reading an
// entire HTTP request, including the body. A zero or negative value means there
// will be no timeout.
//
// Currently used to make operations fail faster in integration tests.
func WithReadTimeout(d time.Duration) func(*Server) {
	return func(s *Server) {
		s.server.ReadTimeout = d
	}
}

// WithWriteTimeout specifies how long the server should wait before timing out
// writing an HTTP response. A zero or negative value means there will be no
// timeout.
//
// Currently used to make operations fail faster in integration tests.
func WithWriteTimeout(d time.Duration) func(*Server) {
	return func(s *Server) {
		s.server.WriteTimeout = d
	}
}

// NewServer creates and returns a new Server backed by the PortService provided.
// It is configured with reasonable defaults, but configuration can be overridden
// using functional options.
func NewServer(addr string, ps PortService, opts ...func(*Server)) *Server {
	mux := &httprouter.Router{
		HandleMethodNotAllowed: false,
		HandleOPTIONS:          false,
	}
	srv := &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		logger: log.New(os.Stderr, "http ", log.LstdFlags),
		Ports:  ps,
	}
	for _, optionFn := range opts {
		optionFn(srv)
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

// ErrorResponse is the HTTP response body delivered for every HTTP request that
// cannot be processed.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ReplyErr examines the error provided and replies to an HTTP request with an
// appropriate HTTP code and payload.
//
// If the error provided is an instance of a ports.Error, the function will use
// the error metadata to construct a meaningful ErrorResponse, and return that
// with an appropriate HTTP status code. If the error is not an instance of
// ports.Error, the function will reply with an HTTP 503 (Service Unavailable)
// and a somewhat generic payload.
//
// The function does not otherwise end the request; the caller should ensure no
// further writes are done to w.
func (s *Server) ReplyErr(w http.ResponseWriter, err error) {
	var (
		statusCode = http.StatusServiceUnavailable // Default HTTP status code.
		errCode    = ports.ErrCodeInternal         // Default response error code.
		errMsg     = "please try again later"      // Default response error message.
	)

	var portsErr *ports.Error
	if ok := errors.As(err, &portsErr); ok {
		errCode, errMsg = portsErr.Code, portsErr.Msg

		// Use the error code to "translate" the error into an appropriate HTTP response.
		switch portsErr.Code {
		case ports.ErrCodeInvalid:
			statusCode = http.StatusBadRequest
		case ports.ErrCodeNotFound:
			statusCode = http.StatusNotFound
		}
	}

	s.Reply(w, statusCode, &ErrorResponse{
		Code:    errCode,
		Message: errMsg,
	})
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
