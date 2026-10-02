package advanced

import (
	"fmt"
	"testing"
)

// ---------------------------------------------------------------------------
// Go Testing, Benchmarking, and Fuzzing
//
// For a Java developer:
// - Go's `testing` package is built into the toolchain (no external JUnit needed).
// - Tests are always in files ending in `_test.go` and in the same package.
// - `TestXxx` are standard tests (like `@Test` in JUnit).
// - `BenchmarkXxx` are for performance testing (built-in JMH-like functionality).
// - `FuzzXxx` are for fuzz testing (available from Go 1.18+, similar to JQF).
// - `ExampleXxx` are executable documentation verified as tests via `// Output:` comments.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Benchmarking (Go's built-in performance measurement)
// Run with: go test -bench=. ./pkg/advanced
// ---------------------------------------------------------------------------

func BenchmarkSum(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(nums)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	var result int
	for i := 0; i < b.N; i++ {
		result = fib(20)
	}
	_ = result
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

// ---------------------------------------------------------------------------
// Fuzz Testing (Go's built-in fuzzing - Go 1.18+)
// Run with: go test -fuzz=FuzzReverse ./pkg/advanced
// ---------------------------------------------------------------------------

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)
		if orig != doubleRev {
			t.Errorf("Before: %q, after double reverse: %q", orig, doubleRev)
		}
	})
}

// Reverse reverses a string by runes (handles multi-byte UTF-8 characters).
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ---------------------------------------------------------------------------
// Example Functions (Documented, executable examples in tests)
// ---------------------------------------------------------------------------

func ExampleSum() {
	nums := []int{1, 2, 3, 4}
	fmt.Println(Sum(nums))
	// Output: 10
}
