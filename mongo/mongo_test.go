package mongo_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/christgf/ports/mongo"
)

func setup(t *testing.T) (*mongo.DB, func()) {
	t.Helper()
	if testing.Short() {
		t.Skipf("skipping test %s in short mode", t.Name())
	}

	const envURI = "PORTS_MONGODB_CONN_URI"
	uri, ok := os.LookupEnv(envURI)
	if !ok {
		t.Skipf("skipping test %s: environment variable %q is not set", t.Name(), envURI)
	}

	db, err := mongo.Open(uri, mongo.WithServerSelectTimeout(time.Second))
	if err != nil {
		t.Fatalf("Open(): %v", err)
	}

	if err := db.Ping(context.TODO()); err != nil {
		t.Errorf("Ping(): %v", err)
	}

	closeFn := func() {
		if err := db.Ports().Drop(context.TODO()); err != nil {
			t.Errorf("Drop(): %v", err)
		}
		if err := db.Close(); err != nil {
			t.Errorf("Close(): %v", err)
		}
	}

	return db, closeFn
}

func TestCreateIndexes(t *testing.T) {
	db, teardown := setup(t)
	t.Cleanup(teardown)

	want := []string{
		"papaya.ports._id_",
		"papaya.ports.id_1",
	}

	got, err := db.CreateIndexes(context.Background())
	if err != nil {
		t.Fatalf("CreateIndexes() returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateIndexes(): index specifications\nhave: %v\nwant: %v", got, want)
	}
}
