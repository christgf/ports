// Package ports defines domain entities and implements business logic.
package ports

import (
	"context"
	"errors"
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
//
// Implementations are expected to return a ports.Error instance with code
// ErrCodeNotFound when a record matching the portID could not be found. In any
// other case, the error will be interpreted as an internal storage error and
// will be handled accordingly.
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
		return &Error{Code: ErrCodeInvalid, Msg: err.Error(), Cause: err}
	}

	if err := s.Ports.InsertPort(ctx, p); err != nil {
		return &Error{Code: ErrCodeInternal, Msg: "could not insert", Cause: err}
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
		return nil, &Error{Code: ErrCodeInvalid, Msg: "port ID should not be empty", Cause: ErrInvalidPortID}
	}

	port, err := s.Ports.FindPort(ctx, portID)
	if err != nil {
		if errors.Is(err, &Error{Code: ErrCodeNotFound}) {
			return nil, err
		}

		return nil, &Error{Code: ErrCodeInternal, Msg: "an unexpected error has occurred", Cause: err}
	}

	return port, nil
}
