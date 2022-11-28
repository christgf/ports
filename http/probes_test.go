package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/christgf/ports"
	"github.com/christgf/ports/http"
)

func TestHandleAlive(t *testing.T) {
	s := http.NewServer(&ports.Service{})

	w := httptest.NewRecorder()
	s.HandleAlive(w, httptest.NewRequest("GET", "/alive", nil))

	if got, want := w.Result().StatusCode, 200; got != want {
		t.Errorf("HandleAlive: have response code %q, want %q", got, want)
	}
}

func TestHandleReady(t *testing.T) {
	s := http.NewServer(&ports.Service{})

	w := httptest.NewRecorder()
	s.HandleReady(w, httptest.NewRequest("GET", "/ready", nil))

	if got, want := w.Result().StatusCode, 200; got != want {
		t.Errorf("HandleReady: have response code %q, want %q", got, want)
	}
}
