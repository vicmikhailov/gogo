// Package errors showcases Go's error handling model.
//
// For a Java developer:
//   - Go does NOT use exceptions (`try-catch-finally`).
//   - Errors are treated as normal values.
//   - Functions return an `error` as the last return value.
//   - Java-ism to avoid: Using `panic()` for control flow or expected errors.
//   - Go way: Only use `panic()` for truly unrecoverable issues (e.g., programmer errors like
//     dereferencing a nil pointer that shouldn't be nil). Otherwise, always return `error`.
package errors

import (
	"errors"
	"fmt"
)

// MyCustomError is a custom error type.
//
// For a Java developer:
//   - Go doesn't have exceptions (`try-catch-finally`).
//   - An "error" is just a value that implements the `error` interface.
//   - The `error` interface only has one method: `Error() string`.
//   - This is like a checked exception, but you handle it as a return value.
//   - Functions return multiple values: `(result, error)`.
//   - `panic` exists (like `RuntimeException`), but it's for unrecoverable errors.
type MyCustomError struct {
	Code    int
	Message string
}

// Error makes MyCustomError satisfy the error interface.
//
// For a Java developer:
// - Equivalent to `public String getMessage()` in Java's `Throwable`.
// - This is the single method required to implement the `error` interface.
func (e *MyCustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// DoSomethingThatFails returns a custom error.
// Java equivalent: `public void doSomething() throws MyCustomException`
func DoSomethingThatFails() error {
	return &MyCustomError{Code: 404, Message: "Resource not found"}
}

// RunErrorsDemo showcases Go's error handling patterns.
//
// For a Java developer:
//   - Error handling is explicit and immediate.
//   - No "bubbling up" unless you explicitly return the error.
//   - `errors.Is` checks for specific error values (like `instanceof` with a singleton).
//   - `errors.As` checks for specific error types (like `instanceof` with a class).
func RunErrorsDemo() {
	fmt.Println("--- Error Handling Demo ---")

	// 1. Basic error handling
	// Java equivalent: `new Exception("a simple error")`
	fmt.Println("1. Basic error:")
	err := errors.New("a simple error")
	if err != nil { // Standard Go idiom: check if err is not nil
		fmt.Println("   Caught error:", err)
	}

	// 2. Custom error types and type assertion
	// Java equivalent: `if (e instanceof MyCustomException) { ... }`
	fmt.Println("2. Custom error and type assertion:")
	err = DoSomethingThatFails()
	if customErr, ok := err.(*MyCustomError); ok {
		fmt.Printf("   Caught custom error: code=%d, message=%s\n", customErr.Code, customErr.Message)
	}

	// 3. errors.Is and errors.As (Go 1.13+)
	// `errors.As` is the modern, preferred way to check for specific error types,
	// especially if errors are "wrapped" (nested).
	fmt.Println("3. errors.As (modern approach):")
	var target *MyCustomError
	if errors.As(err, &target) {
		fmt.Printf("   Successfully retrieved custom error via errors.As: code=%d\n", target.Code)
	}

	// 4. Error Wrapping (%w)
	// Java comparison: `new Exception("context", originalException)`
	// Go idiom: Use `fmt.Errorf` with `%w` to wrap an error while preserving its identity.
	fmt.Println("4. Error wrapping:")
	baseErr := errors.New("database connection failed")
	wrappedErr := fmt.Errorf("failed to get user: %w", baseErr)
	fmt.Printf("   Wrapped error: %v\n", wrappedErr)
	fmt.Printf("   Is base error? %v\n", errors.Is(wrappedErr, baseErr))

	fmt.Println("--- Error Handling Demo End ---")
}
