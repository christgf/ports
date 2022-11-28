package mongo_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/christgf/ports"
)

func TestDBInsertFindPort(t *testing.T) {
	db, teardown := setup(t)
	t.Cleanup(teardown)

	port := ports.Port{
		ID:       "MXACA",
		Name:     "Acapulco",
		Code:     "20101",
		City:     "Acapulco",
		Province: "Guerrero",
		Country:  "Mexico",
		Timezone: "America/Mexico_City",
		UNLocs:   []string{"MXACA"},
		Coords:   []float64{-99.87, 16.85},
	}

	t.Log("FindPort against empty database, we expect an error")
	_, err := db.FindPort(context.TODO(), port.ID)
	if !errors.Is(err, &ports.Error{Code: ports.ErrCodeNotFound}) {
		t.Fatalf("FindPort(%q): have %v, want not found error", port.ID, err)
	}

	t.Logf("Inserting port with port ID %q, expecting no errors", port.ID)
	if err := db.InsertPort(context.Background(), port); err != nil {
		t.Fatalf("InsertPort(): %v", err)
	}

	t.Logf("FindPort with port ID %q, expecting a record and no errors", port.ID)
	p, err := db.FindPort(context.Background(), port.ID)
	if err != nil {
		t.Fatalf("FindPort(): %v", err)
	}

	if got, want := p, &port; !reflect.DeepEqual(got, want) {
		t.Fatalf("FindPort(): port mismatch\nhave: %+v\nwant: %+v\n", got, want)
	}
}

func TestDBInsertPortUpdateExisting(t *testing.T) {
	db, teardown := setup(t)
	t.Cleanup(teardown)

	port := ports.Port{
		ID:       "MXACA",
		Name:     "ACAPULCO",
		City:     "ACAPULCO",
		Province: "GUERRERO",
		Country:  "Mexico",
		Coords:   []float64{-99.87, 16.85},
	}

	t.Logf("Inserting port with port ID %q, expecting no errors", port.ID)
	if err := db.InsertPort(context.Background(), port); err != nil {
		t.Fatalf("InsertPort(): %v", err)
	}

	port.Code = "20101"
	port.Name = "Acapulco"
	port.City = "Acapulco"
	port.Province = "Guerrero"
	port.Timezone = "America/Mexico_City"
	port.UNLocs = []string{"MXACA"}

	t.Logf("Inserting port with port ID %q again with minor differences, expecting no errors", port.ID)
	if err := db.InsertPort(context.Background(), port); err != nil {
		t.Fatalf("InsertPort(): %v", err)
	}

	t.Logf("FindPort with port ID %q, expecting a record and no errors", port.ID)
	p, err := db.FindPort(context.Background(), port.ID)
	if err != nil {
		t.Fatalf("FindPort(): %v", err)
	}

	if got, want := p, &port; !reflect.DeepEqual(got, want) {
		t.Fatalf("FindPort(): port mismatch\nhave: %+v\nwant: %+v\n", got, want)
	}
}
