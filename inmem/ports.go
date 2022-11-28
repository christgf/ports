package inmem

import (
	"context"
	"errors"

	"github.com/christgf/ports"
)

// InsertPort can store ports.Port records in memory.
func (db *DB) InsertPort(_ context.Context, p ports.Port) error {
	db.Lock()
	defer db.Unlock()

	db.data[p.ID] = p

	return nil
}

// ErrNotFound is returned when a ports.Port record could not be found.
var ErrNotFound = errors.New("not found")

// FindPort can retrieve ports.Port records from memory.
func (db *DB) FindPort(_ context.Context, portID string) (*ports.Port, error) {
	db.RLock()
	defer db.RUnlock()

	p, ok := db.data[portID]
	if !ok {
		return nil, ErrNotFound
	}

	return &p, nil
}
