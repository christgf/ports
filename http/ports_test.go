package http_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

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

	srv := http.NewServer(":http", &ports.Service{
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
	}, http.WithReadTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleGetPort(rec, httptest.NewRequest("GET", "/ports?portID=MXACA", nil))

	if got, want := rec.Result().StatusCode, 200; got != want {
		t.Fatalf("HandleGetPort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Fatalf("HandleGetPort: have content type header %q, want %q", got, want)
	}

	wantBody := `{"id":"MXACA","name":"Acapulco","code":"20101","city":"Acapulco","province":"Guerrero","country":"Mexico","timezone":"America/Mexico_City","unlocs":["MXACA"],"coords":[-99.87,16.85]}`
	if gotBody := readAll(t, rec.Result().Body); gotBody != wantBody {
		t.Errorf("HandleGetPort: unexpected response body\nhave: %s\nwant: %s", gotBody, wantBody)
	}
}

func TestHandleGetPortNotFound(t *testing.T) {
	srv := http.NewServer(":http", &ports.Service{
		Ports: &mock.InsertFinder{
			FindPortFn: func(ctx context.Context, portID string) (*ports.Port, error) {
				return nil, &ports.Error{Code: ports.ErrCodeNotFound, Msg: "could not be found"}
			},
		},
	}, http.WithReadTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleGetPort(rec, httptest.NewRequest("GET", "/ports?portID=FOOBAR", nil))

	if got, want := rec.Result().StatusCode, 404; got != want {
		t.Fatalf("HandleGetPort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Fatalf("HandleGetPort: have content type header %q, want %q", got, want)
	}

	wantBody := `{"code":"missing","message":"could not be found"}`
	if gotBody := readAll(t, rec.Result().Body); gotBody != wantBody {
		t.Errorf("HandleGetPort: unexpected response body\nhave: %s\nwant: %s", gotBody, wantBody)
	}
}

func TestHandleGetPortUnexpectedError(t *testing.T) {
	srv := http.NewServer(":http", &ports.Service{
		Ports: &mock.InsertFinder{
			FindPortFn: func(ctx context.Context, portID string) (*ports.Port, error) {
				return nil, errors.New("something went terribly wrong")
			},
		},
	}, http.WithReadTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleGetPort(rec, httptest.NewRequest("GET", "/ports?portID=FOOBAR", nil))

	if got, want := rec.Result().StatusCode, 503; got != want {
		t.Errorf("HandleGetPort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("HandleGetPort: have content type header %q, want %q", got, want)
	}

	wantBody := `{"code":"internal","message":"an unexpected error has occurred"}`
	if gotBody := readAll(t, rec.Result().Body); gotBody != wantBody {
		t.Errorf("HandleGetPort: unexpected response body\nhave: %s\nwant: %s", gotBody, wantBody)
	}
}

func TestHandleStorePort(t *testing.T) {
	srv := http.NewServer(":http", &ports.Service{
		Ports: &mock.InsertFinder{
			InsertPortFn: func(_ context.Context, p ports.Port) error {
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

				if got, want := p, port; !reflect.DeepEqual(got, want) {
					t.Fatalf("InsertPort(): port mismatch\nhave: %+v\nwant: %+v\n", got, want)
				}

				return nil
			},
		},
	}, http.WithWriteTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleStorePort(rec, httptest.NewRequest("GET", "/ports", bytes.NewBufferString(`{
		"ID": "MXACA",
		"Name": "Acapulco",
		"Code": "20101",
		"City": "Acapulco",
		"Province": "Guerrero",
		"Country": "Mexico",
		"Timezone": "America/Mexico_City",
		"UNLocs": ["MXACA"],
		"Coords": [-99.87, 16.85]
	}`)))

	if got, want := rec.Result().StatusCode, 201; got != want {
		t.Fatalf("HandleStorePort: have response code %d, want %d", got, want)
	}
}

func TestHandleStorePortDecodeError(t *testing.T) {
	srv := http.NewServer(":http", &ports.Service{}, http.WithWriteTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleStorePort(rec, httptest.NewRequest("GET", "/ports", bytes.NewBufferString("<xml></xml>")))

	if got, want := rec.Result().StatusCode, 400; got != want {
		t.Errorf("HandleStorePort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("HandleStorePort: have content type header %q, want %q", got, want)
	}

	wantBody := `{"code":"invalid","message":"could not decode"}`
	if gotBody := readAll(t, rec.Result().Body); gotBody != wantBody {
		t.Errorf("HandleStorePort: unexpected response body\nhave: %s\nwant: %s", gotBody, wantBody)
	}
}

func TestHandleStorePortInvalidError(t *testing.T) {
	srv := http.NewServer(":http", &ports.Service{}, http.WithWriteTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleStorePort(rec, httptest.NewRequest("GET", "/ports", bytes.NewBufferString(`{ "Name": "Acapulco" }`)))

	if got, want := rec.Result().StatusCode, 400; got != want {
		t.Errorf("HandleStorePort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("HandleStorePort: have content type header %q, want %q", got, want)
	}

	wantBody := `{"code":"invalid","message":"port ID should not be empty"}`
	if gotBody := readAll(t, rec.Result().Body); gotBody != wantBody {
		t.Errorf("HandleStorePort: unexpected response body\nhave: %s\nwant: %s", gotBody, wantBody)
	}
}

func TestHandleStorePortInsertError(t *testing.T) {
	srv := http.NewServer(":http", &ports.Service{
		Ports: &mock.InsertFinder{
			InsertPortFn: func(_ context.Context, p ports.Port) error {
				return errors.New("the database has gone missing")
			},
		},
	}, http.WithWriteTimeout(time.Second))

	rec := httptest.NewRecorder()
	srv.HandleStorePort(rec, httptest.NewRequest("GET", "/ports", bytes.NewBufferString(`{
		"ID": "MXACA",
		"Name": "Acapulco",
		"Code": "20101",
		"City": "Acapulco",
		"Province": "Guerrero",
		"Country": "Mexico",
		"Timezone": "America/Mexico_City",
		"UNLocs": ["MXACA"],
		"Coords": [-99.87, 16.85]
	}`)))

	if got, want := rec.Result().StatusCode, 503; got != want {
		t.Errorf("HandleStorePort: have response code %d, want %d", got, want)
	}
	if got, want := rec.Result().Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("HandleStorePort: have content type header %q, want %q", got, want)
	}

	wantBody := `{"code":"internal","message":"could not insert"}`
	if gotBody := readAll(t, rec.Result().Body); gotBody != wantBody {
		t.Errorf("HandleStorePort: unexpected response body\nhave: %s\nwant: %s", gotBody, wantBody)
	}
}

func readAll(t *testing.T, src io.ReadCloser) string {
	t.Helper()
	defer func() {
		if err := src.Close(); err != nil {
			t.Logf("io.Closer.Close: %v", err)
		}
	}()

	v, err := io.ReadAll(src)
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}

	return string(bytes.TrimSpace(v))
}
