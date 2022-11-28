// Package main is an HTTP API application.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/christgf/ports"
	"github.com/christgf/ports/http"
	"github.com/christgf/ports/mongo"
)

func main() {
	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancelFn()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		cancelFn()
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	conf := ParseFlags()
	m := Main{
		Conf:   conf,
		Logger: log.New(os.Stdout, "main ", log.LstdFlags),
	}

	if err := m.Run(ctx); err != nil {
		return err
	}

	return nil
}

// Main represents the program, our HTTP API server.
type Main struct {
	Conf   Config
	Logger *log.Logger
}

// Run executes Main, launching our HTTP API. It bootstraps ports.Service with
// the appropriate dependencies, establishes a connection to our MongoDB database
// and configures an HTTP server according to Main.Conf, and begins serving
// incoming HTTP requests. It returns a meaningful error if something goes wrong,
// or nil when the HTTP server is eventually shut down.
func (m Main) Run(ctx context.Context) error {
	mongoDB, err := mongo.Open(m.Conf.MongoDBURI)
	if err != nil {
		return fmt.Errorf("creating MongoDB client: %w", err)
	}

	if err = mongoDB.Ping(ctx); err != nil {
		return fmt.Errorf("pinging MongoDB: %w", err)
	}

	defer func() {
		if err = mongoDB.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error closing MongoDB client: %v", err)
		}
	}()

	if _, err = mongoDB.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("creating MongoDB indexes: %w", err)
	}

	// Create a new ports service with resolved dependencies.
	service := &ports.Service{
		Ports: mongoDB,
	}

	// Set up HTTP server, backed by our ports service implementation.
	server := http.NewServer(m.Conf.HTTPListenAddr, service)

	// Serve.
	m.Logger.Printf("Listening on %s...", m.Conf.HTTPListenAddr)
	if err := server.Serve(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

// Config is the application configuration.
type Config struct {
	HTTPListenAddr string // The listener address for the HTTP server.
	MongoDBURI     string // The MongoDB connection URI.
}

// ParseFlags parses the command line arguments and produces application
// configuration in the form of Config.
//
// It exists as a separate function so that it can be skipped in end-to-end
// tests. Tests can provide their own Config.
func ParseFlags() Config {
	var conf Config
	{
		flag.StringVar(&conf.HTTPListenAddr, "http-listen-addr", getEnvString("PORTS_HTTP_LISTEN_ADDR", ":http"), "HTTP server port")
		flag.StringVar(&conf.MongoDBURI, "mongodb-conn-uri", getEnvString("PORTS_MONGODB_CONN_URI", "mongodb://localhost:27017/ports"), "MongoDB connection URI")
	}
	flag.Parse()

	return conf
}

// getEnvString retrieves the value of the environment variable named by the key.
// If the variable is present in the environment, its value (which may be empty)
// is returned, otherwise fallback is returned.
func getEnvString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}
