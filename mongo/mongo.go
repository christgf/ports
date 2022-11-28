// Package mongo provides MongoDB integration.
package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	defaultConnectTimeout      = 4 * time.Second
	defaultServerSelectTimeout = 8 * time.Second
	defaultCloseTimeout        = 4 * time.Second
)

// DB is attached to a MongoDB database. It maintains a MongoDB connection
// pool and will open and close database connections automatically. It provides
// handles for database collections of interest. It is safe for concurrent use by
// multiple goroutines.
type DB struct {
	*mongo.Database
	client         *mongo.Client
	options        *options.ClientOptions
	closeTimeout   time.Duration
	connectTimeout time.Duration

	// Ports is the ports collection. Declared as a function, so that tests can
	// overwrite the actual collection if they need to.
	Ports func() *mongo.Collection
}

// Names of MongoDB database collections.
const (
	collectionPorts = "ports"
)

// WithServerSelectTimeout specifies how long the driver will wait to find an
// available, suitable server to execute an operation. The default value is
// defaultServerSelectTimeout.
//
// This is currently used to make operations fail faster in tests.
func WithServerSelectTimeout(d time.Duration) func(*DB) {
	return func(db *DB) {
		db.options.SetServerSelectionTimeout(d)
	}
}

// Open will attempt to connect on a specific database on a MongoDB instance. The
// URI parameter should be a valid MongoDB connection URI, with the database name
// being a non-empty string. Note that Open will not return an error if the
// connection cannot be established. To verify the connection, Ping should be
// used in a separate call.
func Open(uri string, opts ...func(*DB)) (*DB, error) {
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, fmt.Errorf("parsing URI: %w", err)
	}

	clientOptions := options.Client().
		SetServerSelectionTimeout(defaultServerSelectTimeout).
		ApplyURI(uri)

	db := &DB{
		options:        clientOptions,
		connectTimeout: defaultConnectTimeout,
		closeTimeout:   defaultCloseTimeout,
	}
	for _, optionFn := range opts {
		optionFn(db)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), db.connectTimeout)
	defer cancelFn()

	client, err := mongo.Connect(ctx, db.options)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	db.client = client
	db.Database = client.Database(cs.Database)

	// Keep a reference of database collections.
	db.Ports = func() *mongo.Collection {
		return db.Collection(collectionPorts)
	}

	return db, nil
}

// CreateIndexes will create/update indexes for database collections. It returns
// collection indexes, or an error if an index operation fails. It will return
// early in case of error, potentially without running all index operations.
//
// Note that creating MongoDB indexes is an idempotent operation, so
// the routine should only create indexes that don't already exist.
func (db *DB) CreateIndexes(ctx context.Context) ([]string, error) {
	// Ports ID index, each port should have a unique identifier.
	const portIDIndex = "id_1"
	if _, err := db.Ports().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
		},
		Options: options.Index().SetName(portIDIndex).SetUnique(true),
	}); err != nil {
		return nil, fmt.Errorf("creating index %q: %w", portIDIndex, err)
	}

	// Retrieve index specifications.
	specs, err := db.Ports().Indexes().ListSpecifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving index specs: %v", err)
	}
	indexes := make([]string, len(specs))
	for i, spec := range specs {
		indexes[i] = fmt.Sprintf("%s.%s", spec.Namespace, spec.Name)
	}

	return indexes, nil
}

// Ping sends a ping command to verify that the client is connected to the
// database. It should be used when the application starts to verify the database
// connection.
func (db *DB) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}

// Close will disconnect the MongoDB connection handle. It will wait for a
// reasonable period of time for in-use connections to be closed and to be
// returned to the connection pool.
func (db *DB) Close() error {
	ctx, cancelFn := context.WithTimeout(context.Background(), db.closeTimeout)
	defer cancelFn() // connections close immediately if context expires.
	return db.client.Disconnect(ctx)
}
