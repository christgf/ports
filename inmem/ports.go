package inmem

import (
	"context"

	"github.com/christgf/ports"
)

// InsertPort can store ports.Port records in memory.
func (db *DB) InsertPort(_ context.Context, p ports.Port) error {
	db.Lock()
	defer db.Unlock()

	db.data[p.ID] = p

	return nil
}

// FindPort can retrieve ports.Port records from memory.
func (db *DB) FindPort(_ context.Context, portID string) (*ports.Port, error) {
	db.RLock()
	defer db.RUnlock()

	p, ok := db.data[portID]
	if !ok {
		return nil, &ports.Error{Code: ports.ErrCodeNotFound, Msg: "port not found"}
	}

	return &p, nil
}
