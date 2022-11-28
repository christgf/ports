// Package ports defines domain entities and implements business logic.
package ports

import (
	"context"
	"errors"
	"fmt"
)

// Port represents a port.
type Port struct {
	ID       string // Unique identifier for Port.
	Name     string
	Code     string
	City     string
	Province string
	Country  string
	Alias    []string
	Regions  []string
	Timezone string
	UNLocs   []string
	Coords   []float64
}

// Errors for unexpected or unsupported values for Port fields.
var (
	ErrInvalidPortID   = errors.New("port ID should not be empty")
	ErrInvalidPortName = errors.New("port name should not be empty")
	ErrInvalidPortCode = errors.New("port code should not be empty")
)

// Validate examines Port fields and returns an appropriate error if any of the
// fields hold unexpected or unsupported values.
func Validate(p Port) error {
	switch {
	case p.ID == "":
		return ErrInvalidPortID
	case p.Name == "":
		return ErrInvalidPortName
	case p.Code == "":
		return ErrInvalidPortCode
	default:
		return nil
	}
}

// Inserter can insert Port records in storage.
type Inserter interface {
	InsertPort(ctx context.Context, p Port) error
}

// Finder can retrieve Port records from storage.
type Finder interface {
	FindPort(ctx context.Context, portID string) (*Port, error)
}

// InsertFinder groups Inserter and Finder capabilities for Port records.
type InsertFinder interface {
	Inserter
	Finder
}

// Service manages Port instances and records.
type Service struct {
	Ports InsertFinder // Port record storage.
}

// StorePort records port information in storage. It returns an error if the
// information provided is unexpected or invalid, if the underlying storage
// system fails, or if the context is cancelled before the operation is
// completed. It returns nothing if the operation is successful.
func (s *Service) StorePort(ctx context.Context, p Port) error {
	if err := Validate(p); err != nil {
		return fmt.Errorf("validating: %w", err)
	}

	if err := s.Ports.InsertPort(ctx, p); err != nil {
		return fmt.Errorf("inserting: %w", err)
	}

	return nil
}

// GetPortByID retrieves port information from storage, based on the port
// identifier provided. It returns an appropriate error if the underlying storage
// system fails, if a record matching the identifier is not found, or if the
// context is cancelled before the operation is completed. The portID argument
// should not be empty.
func (s *Service) GetPortByID(ctx context.Context, portID string) (*Port, error) {
	if portID == "" {
		return nil, ErrInvalidPortID
	}

	port, err := s.Ports.FindPort(ctx, portID)
	if err != nil {
		return nil, fmt.Errorf("finding: %w", err)
	}

	return port, nil
}
