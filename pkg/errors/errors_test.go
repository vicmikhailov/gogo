package errors

import (
	"errors"
	"testing"
)

/**
 * ===========================================================================
 * Error Handling Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Use `errors.As` to check if an error is of a specific type (like `instanceof`).
 * - Use `errors.Is` to check for specific sentinel error values.
 */

func TestCustomError(t *testing.T) {
	err := DoSomethingThatFails()

	// 1. Check for specific type using errors.As
	var target *MyCustomError
	if !errors.As(err, &target) {
		t.Errorf("Expected MyCustomError, got %T", err)
	} else {
		if target.Code != 404 {
			t.Errorf("Expected code 404, got %d", target.Code)
		}
	}

	// 2. String representation
	expected := "Error 404: Resource not found"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}
