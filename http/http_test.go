package http_test

import (
	"context"
	"testing"

	"github.com/christgf/ports"
	"github.com/christgf/ports/http"
)

func TestServe(t *testing.T) {
	s := http.NewServer(&ports.Service{})

	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	if err := s.Serve(ctx); err != nil {
		t.Errorf("Serve: cancelled context should gracefully shutdown, got %v", err)
	}
}
