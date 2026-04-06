package advanced // Declaring the package as advanced for testing and benchmarking.

import ( // Starting the import block.
	"fmt"     // Importing fmt for example function output.
	"testing" // Importing the testing package for tests, benchmarks, and fuzzing.
) // Closing the import block.

/**
 * ===========================================================================
 * Go Testing, Benchmarking, and Fuzzing
 * ===========================================================================
 *
 * For a Java developer:
 * - Go's `testing` package is built into the toolchain (no external JUnit needed).
 * - Tests are always in files ending in `_test.go` and in the same package.
 * - `TestXxx` are standard tests (like `@Test` in JUnit).
 * - `BenchmarkXxx` are for performance testing (built-in JMH-like functionality).
 * - `FuzzXxx` are for fuzz testing (available from Go 1.18+, similar to JQF).
 * - `ExampleXxx` are executable documentation that can be verified as tests.
 */ // Block comment explaining Go testing features to Java developers.

// ---------------------------------------------------------------------------
// Benchmarking (Go's built-in performance measurement) // Section header for Benchmarking.
// ---------------------------------------------------------------------------

// To run: go test -bench=. ./pkg/demo // Instructions on how to run benchmarks.
func BenchmarkSum(b *testing.B) { // Benchmark function for the Sum function.
	nums := make([]int, 1000) // Creating a slice of 1000 integers.
	for i := range nums {     // Iterating through the slice to initialize it.
		nums[i] = i // Setting each element to its index.
	} // End of initialization loop.
	b.ResetTimer()             // Resetting the benchmark timer to exclude setup time.
	for i := 0; i < b.N; i++ { // Running the benchmark b.N times.
		Sum(nums) // Calling the Sum function.
	} // End of benchmark loop.
} // Closing the BenchmarkSum function.

func BenchmarkFibonacci(b *testing.B) { // Benchmark function for the Fibonacci calculation.
	var result int             // Variable to store the result and prevent optimization.
	for i := 0; i < b.N; i++ { // Running the benchmark b.N times.
		result = fib(20) // Calculating the 20th Fibonacci number.
	} // End of benchmark loop.
	_ = result // Using the result to avoid compiler optimization.
} // Closing the BenchmarkFibonacci function.

func fib(n int) int { // Recursive Fibonacci function.
	if n <= 1 { // Base case: if n is 0 or 1.
		return n // Return n.
	} // End of base case.
	return fib(n-1) + fib(n-2) // Recursive step: sum of the two preceding numbers.
} // Closing the fib function.

// ---------------------------------------------------------------------------
// Fuzz Testing (Go's built-in fuzzing - Go 1.18+) // Section header for Fuzz Testing.
// ---------------------------------------------------------------------------

// To run: go test -fuzz=FuzzReverse ./pkg/demo // Instructions on how to run fuzz tests.
func FuzzReverse(f *testing.F) { // Fuzz test function for the Reverse function.
	testcases := []string{"Hello, world", " ", "!12345"} // Defining initial seed corpus.
	for _, tc := range testcases {                       // Iterating through seed test cases.
		f.Add(tc) // Providing the seed corpus to the fuzzer.
	} // End of corpus addition loop.
	f.Fuzz(func(t *testing.T, orig string) { // Defining the fuzz target function.
		rev := Reverse(orig)      // Reversing the original string.
		doubleRev := Reverse(rev) // Reversing the reversed string.
		if orig != doubleRev {    // Checking if double reversal returns the original string.
			t.Errorf("Before: %q, after double reverse: %q", orig, doubleRev) // Logging error if mismatch.
		} // End of property check.
		// Check for valid UTF-8 // Note on further possible checks.
		// (Optional, but good for strings) // Note on optional checks.
	}) // End of f.Fuzz call.
} // Closing the FuzzReverse function.

// Reverse reverses a string by runes (handles multi-byte characters). // Comment for the Reverse function.
func Reverse(s string) string { // Reverse function taking a string and returning a string.
	runes := []rune(s)                                    // Converting the string to a slice of runes (UTF-8 points).
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { // Using two pointers to reverse the slice in-place.
		runes[i], runes[j] = runes[j], runes[i] // Swapping runes at positions i and j.
	} // End of reversal loop.
	return string(runes) // Converting the slice of runes back to a string and returning it.
} // Closing the Reverse function.

// ---------------------------------------------------------------------------
// Example Functions (Documented examples in tests) // Section header for Example functions.
// ---------------------------------------------------------------------------

func ExampleSum() { // Example function for the Sum function (verifiable documentation).
	nums := []int{1, 2, 3, 4} // Initializing a slice of integers.
	fmt.Println(Sum(nums))    // Printing the sum of the integers.
	// Output: 10
} // Closing the ExampleSum function.
