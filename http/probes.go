package http

import "net/http"

// HandleAlive handles liveliness probes. It simply returns OK for now.
func (s *Server) HandleAlive(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// HandleReady handles the service health probe. It is quite optimistic and
// always returns OK without actually checking if the service and/or its
// components are healthy and ready to process requests.
func (s *Server) HandleReady(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
