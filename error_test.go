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
			res: "(missing) resource missing",
		},
	}

	for _, tt := range tests {
		if got, want := tt.err.Error(), tt.res; got != want {
			t.Errorf("err.Error(): have %q, want %q", got, want)
		}
	}
}

func TestErrorIs(t *testing.T) {
	var err error = &ports.Error{Code: ports.ErrCodeInvalid}
	if !errors.Is(err, &ports.Error{Code: ports.ErrCodeInvalid}) {
		t.Errorf("expecting Error with CodeInvalid, got: %s", err)
	}

	target := errors.New("something entirely different")
	if errors.Is(err, target) {
		t.Errorf("Is(%#v, %#v) should not be true", err, target)
	}
}
