// Package inmem provides in-memory implementations of various components.
package inmem

import (
	"sync"

	"github.com/christgf/ports"
)

// DB is an in-memory implementation of ports.InsertFinder.
type DB struct {
	sync.RWMutex
	data map[string]ports.Port
}

// Open instantiates and returns a new DB.
func Open() *DB {
	return &DB{
		data: make(map[string]ports.Port, 0),
	}
}
