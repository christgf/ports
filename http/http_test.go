package http_test

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/christgf/ports/http"
)

func TestServeError(t *testing.T) {
	srv := http.NewServer("$", nil, http.WithLoggerOutput(io.Discard))

	err := srv.Serve(context.TODO())
	if err == nil {
		t.Fatalf("Serve: expected error, got nothing")
	}

	var netOpErr *net.OpError
	if !errors.As(err, &netOpErr) {
		t.Errorf("Serve: expected %T but got %T %q", netOpErr, err, err)
	}
}

func TestServeShutdown(t *testing.T) {
	srv := http.NewServer(":http", nil, http.WithLoggerOutput(io.Discard))
	ctx, cancelFn := context.WithCancel(context.Background())
	const cancelAfter = 200 * time.Millisecond

	var wg sync.WaitGroup

	t.Logf("Starting server on %q, begin serving...", srv.Addr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(ctx); err != nil {
			t.Errorf("Serve: expected no error, got %v", err)
		}
	}()

	t.Logf("Cancelling Serve() context after %v, expecting server to shut down gracefully...", cancelAfter)
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.AfterFunc(cancelAfter, cancelFn)
	}()

	wg.Wait()
}
