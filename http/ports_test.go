package http_test

import (
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/christgf/ports"
	"github.com/christgf/ports/http"
	"github.com/christgf/ports/mock"
)

func TestHandleGetPort(t *testing.T) {
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

	s := http.NewServer(&ports.Service{
		Ports: &mock.InsertFinder{
			FindPortFn: func(ctx context.Context, portID string) (*ports.Port, error) {
				if got, want := portID, port.ID; got != want {
					t.Fatalf("FindPort(): have port ID %q, want %q", got, want)
				}

				return &ports.Port{
					ID:       "MXACA",
					Name:     "Acapulco",
					Code:     "20101",
					City:     "Acapulco",
					Country:  "Mexico",
					Province: "Guerrero",
					Timezone: "America/Mexico_City",
					UNLocs:   []string{"MXACA"},
					Coords:   []float64{-99.87, 16.85},
				}, nil
			},
		},
	})

	rec := httptest.NewRecorder()
	s.HandleGetPort(rec, httptest.NewRequest("GET", "/ports?portID=MXACA", nil))
	t.Cleanup(func() {
		_ = rec.Result().Body.Close()
	})

	if got, want := rec.Result().StatusCode, 200; got != want {
		t.Fatalf("HandleGetPort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Fatalf("HandleGetPort: have content type header %q, want %q", got, want)
	}

	resBody, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatalf("HandleGetPort: reading response body: %v", err)
	}

	wantBody := `{"ID":"MXACA","Name":"Acapulco","Code":"20101","City":"Acapulco","Province":"Guerrero","Country":"Mexico","Alias":null,"Regions":null,"Timezone":"America/Mexico_City","UNLocs":["MXACA"],"Coords":[-99.87,16.85]}`
	if gotBody := string(bytes.TrimSpace(resBody)); gotBody != wantBody {
		t.Errorf("HandleGetPort: unexpected response body\n%s", gotBody)
	}
}

func TestHandleStorePort(t *testing.T) {
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

	s := http.NewServer(&ports.Service{
		Ports: &mock.InsertFinder{
			InsertPortFn: func(_ context.Context, p ports.Port) error {
				if got, want := p, port; !reflect.DeepEqual(got, want) {
					t.Fatalf("InsertPort(): port mismatch\nhave: %+v\nwant: %+v\n", got, want)
				}

				return nil
			},
		},
	})

	reqBody := bytes.NewBufferString(`{
		"ID": "MXACA",
		"Name": "Acapulco",
		"Code": "20101",
		"City": "Acapulco",
		"Province": "Guerrero",
		"Country": "Mexico",
		"Timezone": "America/Mexico_City",
		"UNLocs": ["MXACA"],
		"Coords": [-99.87, 16.85]
	}`)
	rec := httptest.NewRecorder()
	s.HandleStorePort(rec, httptest.NewRequest("GET", "/ports", reqBody))
	t.Cleanup(func() {
		_ = rec.Result().Body.Close()
	})

	if got, want := rec.Result().StatusCode, 201; got != want {
		t.Fatalf("HandleStorePort: have response code %d, want %d", got, want)
	}
}
