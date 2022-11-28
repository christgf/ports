package inmem_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/christgf/ports"
	"github.com/christgf/ports/inmem"
)

func TestDBInsertFindPort(t *testing.T) {
	port := ports.Port{
		ID:       "MXACA",
		Name:     "Acapulco",
		Code:     "20101",
		City:     "Acapulco",
		Country:  "Mexico",
		Province: "Guerrero",
		Timezone: "America/Mexico_City",
		UNLocs:   []string{"MXACA"},
		Coords:   []float64{-99.87, 16.85},
	}

	db := inmem.Open()

	t.Logf("FindPort against empty database, we expect %q error", inmem.ErrNotFound)
	_, err := db.FindPort(context.TODO(), port.ID)
	if got, want := err, inmem.ErrNotFound; !errors.Is(got, want) {
		t.Fatalf("FindPort(%q): have %v, want %v", port.ID, got, want)
	}

	t.Logf("Inserting port with port ID %q, expecting no errors", port.ID)
	if err := db.InsertPort(context.TODO(), port); err != nil {
		t.Fatalf("InsertPort(%+v): %v", port.ID, err)
	}

	t.Logf("FindPort with port ID %q, expecting a record and no errors", port.ID)
	p, err := db.FindPort(context.TODO(), port.ID)
	if err != nil {
		t.Fatalf("FindPort(%q): %v", port.ID, err)
	}

	if got, want := p, &port; !reflect.DeepEqual(got, want) {
		t.Fatalf("FindPort(): port mismatch\nhave: %+v\nwant: %+v\n", got, want)
	}
}
