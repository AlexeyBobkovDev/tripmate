package core_errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	cause := errors.New("test error")

	err := NewError(cause, CategoryBadRequest)

	if err.Error() != "test error" {
		t.Fatalf("Error() = %q, want %q", err.Error(), "test error")
	}

	if !errors.Is(err, cause) {
		t.Fatal("errors.Is() = false, want true")
	}

	codedErr, ok := err.(*CodedError)
	if !ok {
		t.Fatal("error is not *CodedError")
	}

	if codedErr.Category != CategoryBadRequest {
		t.Fatalf(
			"Category = %v, want %v",
			codedErr.Category,
			CategoryBadRequest,
		)
	}
}

func TestNewErrorWrapped(t *testing.T) {
	err := fmt.Errorf(
		"create user: %w",
		ErrInvalidUserEmail,
	)

	var codedErr *CodedError

	if !errors.As(err, &codedErr) {
		t.Fatal("errors.As() = false, want true")
	}

	if codedErr.Category != CategoryBadRequest {
		t.Fatalf(
			"Category = %v, want %v",
			codedErr.Category,
			CategoryBadRequest,
		)
	}
}
