package iosystem

import (
	"os"
	"testing"
)

/**
 * ===========================================================================
 * I/O System Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Go's `os` and `io` packages are equivalent to `java.io` and `java.nio`.
 * - `defer` is excellent for cleanup (like deleting temporary files).
 */

func TestRunIOSystemDemo(t *testing.T) {
	// We just want to make sure it runs without crashing.
	// We redirect stdout if we wanted to check output, but for now,
	// a simple call is enough to ensure no nil pointer dereferences etc.

	// Ensure we don't have leftover files
	defer os.Remove("test_demo.txt")

	RunIOSystemDemo()
}
