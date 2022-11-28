// Package mock defines mocks for testing.
//
// Mock functions should not return errors if they are not set, allowing the zero
// value to be used as a valid mock.
package mock

import (
	"context"
	"sync"

	"github.com/christgf/ports"
)

// InsertFinder is a mock implementation of ports.InsertFinder.
type InsertFinder struct {
	InsertPortFn func(ctx context.Context, p ports.Port) error
	FindPortFn   func(ctx context.Context, portID string) (*ports.Port, error)

	sync.Mutex
	InsertPortCalls int
	FindPortCalls   int
}

// InsertPort invokes the mock implementation.
func (m *InsertFinder) InsertPort(ctx context.Context, p ports.Port) error {
	m.Lock()
	m.InsertPortCalls++
	m.Unlock()

	if m.InsertPortFn == nil {
		return nil
	}

	return m.InsertPortFn(ctx, p)
}

// FindPort invokes the mock implementation.
func (m *InsertFinder) FindPort(ctx context.Context, portID string) (*ports.Port, error) {
	m.Lock()
	m.FindPortCalls++
	m.Unlock()

	if m.FindPortFn == nil {
		return &ports.Port{}, nil
	}

	return m.FindPortFn(ctx, portID)
}
