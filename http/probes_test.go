package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/christgf/ports"
	"github.com/christgf/ports/http"
)

func TestHandleAlive(t *testing.T) {
	s := http.NewServer(":http", &ports.Service{})

	w := httptest.NewRecorder()
	s.HandleAlive(w, httptest.NewRequest("GET", "/alive", nil))

	if got, want := w.Result().StatusCode, 200; got != want {
		t.Errorf("HandleAlive(): have response status code %d, want %d", got, want)
	}
}

func TestHandleReady(t *testing.T) {
	s := http.NewServer(":http", &ports.Service{})

	w := httptest.NewRecorder()
	s.HandleReady(w, httptest.NewRequest("GET", "/ready", nil))

	if got, want := w.Result().StatusCode, 200; got != want {
		t.Errorf("HandleReady(): have response status code %d, want %d", got, want)
	}
}
