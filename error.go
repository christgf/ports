package ports

import "fmt"

// Machine-parseable error codes used for error translation and for automating
// behavior if and where possible.
const (
	ErrCodeInternal = "internal" // Unexpected failure.
	ErrCodeInvalid  = "invalid"  // Invalid arguments or input.
	ErrCodeNotFound = "missing"  // Requested entity or record not found.
)

// Error is a ports service error.
type Error struct {
	Code  string // Machine-parseable error code.
	Msg   string // Human-readable error message.
	Cause error  // The underlying cause of this error, if any.
}

// Unwrap the underlying error when using the %w verb.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Error implements the built-in error interface.
func (e *Error) Error() string {
	msg := e.Msg
	if len(msg) == 0 {
		if e.Cause != nil {
			msg = e.Cause.Error()
		}
	}

	return fmt.Sprintf("(%s) %s", e.Code, msg)
}

// Is allows to check if a service error has a specific code.
//
//	if errors.Is(err, &Error{Code: ErrCodeInvalid}) {
//		// Do something if err has code invalid.
//	}
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.Code == t.Code
}
