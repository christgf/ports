package ports_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/christgf/ports"
	"github.com/christgf/ports/mock"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		port ports.Port
		err  error
	}{
		{
			port: ports.Port{},
			err:  ports.ErrInvalidPortID,
		},
		{
			port: ports.Port{ID: "AEAJM"},
			err:  ports.ErrInvalidPortName,
		},
		{
			port: ports.Port{ID: "AEAJM", Name: "Ajman"},
			err:  ports.ErrInvalidPortCode,
		},
		{
			port: ports.Port{ID: "AEAJM", Name: "Ajman", Code: "52000"},
			err:  nil,
		},
	}

	for _, tt := range tests {
		gotErr := ports.Validate(tt.port)
		if wantErr := tt.err; !errors.Is(gotErr, wantErr) {
			t.Errorf("Validate(%+v): have %v, want %v", tt.port, gotErr, wantErr)
		}
	}
}

func TestServiceStorePortValidateError(t *testing.T) {
	s := &ports.Service{}

	err := s.StorePort(context.TODO(), ports.Port{})
	if err == nil {
		t.Fatal("StorePort(): expected validation error, got nothing")
	}

	if gotErr, wantErr := err, ports.ErrInvalidPortID; !errors.Is(gotErr, wantErr) {
		t.Errorf("StorePort(): have %q, want port ID validation error", err)
	}
}

func TestServiceStorePortInsertError(t *testing.T) {
	port := ports.Port{
		ID:   "MXACA",
		Name: "Acapulco",
		Code: "20101",
	}
	wantErr := errors.New("something went wrong")

	s := &ports.Service{
		Ports: &mock.InsertFinder{
			InsertPortFn: func(_ context.Context, _ ports.Port) error {
				return wantErr
			},
		},
	}

	if err := s.StorePort(context.TODO(), port); !errors.Is(err, wantErr) {
		t.Errorf("StorePort(): have %v, want %v", err, wantErr)
	}
}

func TestServiceStorePortOK(t *testing.T) {
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

	s := &ports.Service{
		Ports: &mock.InsertFinder{
			InsertPortFn: func(_ context.Context, p ports.Port) error {
				if got, want := p, port; !reflect.DeepEqual(got, want) {
					t.Fatalf("InsertPort(): port mismatch\nhave: %+v\nwant: %+v\n", got, want)
				}

				return nil
			},
		},
	}

	if err := s.StorePort(context.TODO(), port); err != nil {
		t.Errorf("StorePort(): %v", err)
	}
}
