// Package main is an HTTP API application.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/christgf/ports"
	"github.com/christgf/ports/http"
	"github.com/christgf/ports/inmem"
)

func main() {
	ctx := context.Background() // todo(christgf): cancel on signal

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	db := inmem.Open()
	service := &ports.Service{
		Ports: db,
	}

	server := http.NewServer(service)
	if err := server.Serve(ctx); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}
