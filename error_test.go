package ports_test

import (
	"errors"
	"testing"

	"github.com/christgf/ports"
)

func TestErrorError(t *testing.T) {
	tests := []struct {
		err *ports.Error
		res string
	}{
		{
			err: &ports.Error{Code: ports.ErrCodeInternal, Cause: errors.New("unexpected failure")},
			res: "(internal) unexpected failure",
		},
		{
			err: &ports.Error{Code: ports.ErrCodeInvalid, Msg: "something is invalid"},
			res: "(invalid) something is invalid",
		},
		{
			err: &ports.Error{Code: ports.ErrCodeNotFound, Msg: "could not be found", Cause: errors.New("resource missing")},
			res: "(missing) could not be found",
		},
	}

	for _, tt := range tests {
		if got, want := tt.err.Error(), tt.res; got != want {
			t.Errorf("err.Error(): have %q, want %q", got, want)
		}
	}
}

func TestErrorIs(t *testing.T) {
	const code = "irrelevant"

	var err error = &ports.Error{Code: code}
	if !errors.Is(err, &ports.Error{Code: code}) {
		t.Errorf("Is(%T) expected true for error code %q, got false", err, code)
	}

	target := errors.New("something entirely different")
	if errors.Is(err, target) {
		t.Errorf("Is(%#v, %#v) should not be true", err, target)
	}
}
