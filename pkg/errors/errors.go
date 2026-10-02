// Package errors showcases Go's error handling model.
//
// For a Java developer:
//   - Go does NOT use exceptions (`try-catch-finally`).
//   - Errors are treated as normal values.
//   - Functions return an `error` as the last return value: `(result, error)`.
//   - Java-ism to avoid: Using `panic()` for control flow or expected errors.
//   - Go way: Only use `panic()` for truly unrecoverable issues (e.g. nil dereferences or programmer bugs).
//     Otherwise, always return `error`.
package errors

import (
	"errors"
	"fmt"
)

// MyCustomError is a custom error type implementing the `error` interface.
//
// For a Java developer:
//   - Go doesn't have exceptions (`try-catch-finally`).
//   - An "error" is just any value that implements the `error` interface.
//   - The `error` interface only requires a single method: `Error() string`.
//   - Similar to a checked exception in Java, but explicitly returned as a value.
type MyCustomError struct {
	Code    int
	Message string
}

// Error makes MyCustomError satisfy the `error` interface.
// Java equivalent: `public String getMessage()` on `Throwable`.
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
//   - Error handling is explicit and immediate (no silent stack unwinding).
//   - `errors.Is` checks for specific sentinel error values (like `==` or singleton check).
//   - `errors.As` inspects error chains for specific error types (like `instanceof`).
func RunErrorsDemo() {
	fmt.Println("--- Error Handling Demo ---")

	// 1. Basic error handling
	// Java equivalent: `new Exception("a simple error")`
	// Standard Go idiom: check `if err != nil`
	fmt.Println("1. Basic error:")
	err := errors.New("a simple error")
	if err != nil {
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
	// `errors.As` is the preferred way to find a target error type in an error chain.
	fmt.Println("3. errors.As (modern approach):")
	var target *MyCustomError
	if errors.As(err, &target) {
		fmt.Printf("   Successfully retrieved custom error via errors.As: code=%d\n", target.Code)
	}

	// 4. Error Wrapping (%w)
	// Java comparison: `new Exception("context", originalException)`
	// Go idiom: Use `fmt.Errorf` with `%w` to wrap an error while preserving its identity for `errors.Is`.
	fmt.Println("4. Error wrapping:")
	baseErr := errors.New("database connection failed")
	wrappedErr := fmt.Errorf("failed to get user: %w", baseErr)
	fmt.Printf("   Wrapped error: %v\n", wrappedErr)
	fmt.Printf("   Is base error? %v\n", errors.Is(wrappedErr, baseErr))

	fmt.Println("--- Error Handling Demo End ---")
}
