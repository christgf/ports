// Package main is a command-line utility for importing port records from a
// JSON file into a database.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/christgf/ports"
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
	var m Main
	{
		flag.StringVar(&m.FilePath, "f", "testdata/ports.json", "Path to JSON file")
	}
	flag.Parse()

	m.Logger = log.New(os.Stdout, "main ", log.LstdFlags)

	if err := m.Run(ctx); err != nil {
		return err
	}

	return nil
}

// Main represents the program, our HTTP API server.
type Main struct {
	FilePath string
	Logger   *log.Logger
}

// Run executes Main. It will attempt to open the file defined by Main.FilePath
// for reading, decode its contents into ports.Port structs using input
// streaming, printing ports information to os.Stdout in the process. The file
// should contain ports information in JSON format.
//
// The format of the file should be one big JSON object, containing port
// information described by port identifiers as object fields. Example:
//
//	{
//	 "AEAJM": {
//	   "name": "Ajman",
//	   "city": "Ajman",
//	   "country": "United Arab Emirates",
//	   "coordinates": [
//	     55.5136433,
//	     25.4052165
//	   ],
//	 },
//	 "AEAUH": {
//	   "name": "Abu Dhabi",
//	   "coordinates": [
//	     54.37,
//	     24.47
//	   ],
//	   "city": "Abu Dhabi"
//	 },
//	 ...
//
// The JSON decoder used expects the file to be in this exact format, and should
// fail in any other case. Failures are returned as meaningful errors.
//
// The file is closed before the function is returned.
func (m Main) Run(ctx context.Context) error {

	f, err := os.Open(m.FilePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			m.Logger.Printf("closing file: %v", err)
		}
	}()

	decoder := json.NewDecoder(f)

	// Read first, opening token, `[` or `{`
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("decoding opening token: %w", err)
	}

	var i int
	var p ports.Port
	for decoder.More() {
		// Check for context cancellation, abort if context is cancelled.
		if err := ctx.Err(); err != nil {
			return err
		}

		i++
		portID, err := decoder.Token() // Decode the port unique identifier for the port.
		if err != nil {
			return fmt.Errorf("decoding port ID: %w", err)
		}

		// Decode the rest of the information.
		if err := decoder.Decode(&p); err != nil {
			return fmt.Errorf("decoding port: %v", err)
		}

		// Assign the port identifier to the ports.Port.
		p.ID = fmt.Sprintf("%s", portID)

		// Log and proceed.
		m.Logger.Printf("%d: Port: %v", i, p)
	}

	// Read last, closing token.
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("decoding closing token: %w", err)
	}

	return nil
}
