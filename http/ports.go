package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/christgf/ports"
)

// port is the representation of ports.Port as a JSON document.
type port struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	City     string    `json:"city"`
	Province string    `json:"province"`
	Country  string    `json:"country"`
	Alias    []string  `json:"alias,omitempty"`
	Regions  []string  `json:"regions,omitempty"`
	Timezone string    `json:"timezone,omitempty"`
	UNLocs   []string  `json:"UNLocs,omitempty"`
	Coords   []float64 `json:"coords,omitempty"`
}

// HandleGetPort todo
func (s *Server) HandleGetPort(w http.ResponseWriter, r *http.Request) {
	portID := r.URL.Query().Get("portID")

	p, err := s.Ports.GetPortByID(r.Context(), portID)
	if err != nil {
		s.ReplyErr(w, err)
		return
	}

	s.Reply(w, http.StatusOK, p)
}

// HandleStorePort todo
func (s *Server) HandleStorePort(w http.ResponseWriter, r *http.Request) {
	var p port
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		s.ReplyErr(w, errors.New("decode error"))
		return
	}

	if err := s.Ports.StorePort(r.Context(), ports.Port{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		City:     p.City,
		Province: p.Province,
		Country:  p.Country,
		Alias:    p.Alias,
		Regions:  p.Regions,
		Timezone: p.Timezone,
		UNLocs:   p.UNLocs,
		Coords:   p.Coords,
	}); err != nil {
		s.ReplyErr(w, err)
		return
	}

	s.Reply(w, http.StatusCreated, nil)
}
