package basictypes // Declaring the package as basictypes for testing.

import ( // Starting the import block for test dependencies.
	"encoding/json" // Importing json for testing JSON marshalling and unmarshalling.
	"reflect"       // Importing reflect for structural equality checks with DeepEqual.
	"strings"       // Importing strings for string content verification.
	"testing"       // Importing the testing package for Go's unit testing framework.
) // Closing the import block.

/**
 * ===========================================================================
 * Basic Types & Standard Library Testing (For a Java developer)
 * ===========================================================================
 *
 * In Go, testing slices and maps often uses `reflect.DeepEqual` for structural
 * comparison, as the `==` operator is not defined for slices and only for
 * simple maps (comparing to nil).
 *
 * Table-driven tests are used here to test string manipulation and JSON logic.
 */ // Block comment explaining testing concepts for basic types.

func TestReverseString(t *testing.T) { // Test function for the ReverseString utility.
	tests := []struct { // Defining a table of test cases.
		name     string // Name of the test case.
		input    string // Input string for the function.
		expected string // Expected reversed output.
	}{ // Starting the test cases.
		{"Empty string", "", ""},                    // Case: empty string.
		{"Single character", "a", "a"},              // Case: single character.
		{"Simple word", "hello", "olleh"},           // Case: standard word.
		{"Palindrome", "racecar", "racecar"},        // Case: palindrome word.
		{"Multi-byte characters", "こんにちは", "はちにんこ"}, // Case: UTF-8 characters.
		{"Mixed", "Go 1.23", "32.1 oG"},             // Case: mixed alphanumeric.
	} // End of test cases.

	for _, tt := range tests { // Iterating through each test case in the table.
		t.Run(tt.name, func(t *testing.T) { // Running a subtest for each case.
			result := ReverseString(tt.input) // Calling the function being tested.
			if result != tt.expected {        // Checking if the result matches the expectation.
				t.Errorf("ReverseString(%q) = %q; want %q", tt.input, result, tt.expected) // Logging error on mismatch.
			} // End of result check.
		}) // End of subtest.
	} // End of test loop.
} // Closing the TestReverseString function.

func TestIsPalindrome(t *testing.T) { // Test function for the IsPalindrome utility.
	tests := []struct { // Defining a table of test cases.
		name     string // Name of the test case.
		input    string // Input string for the function.
		expected bool   // Expected boolean result.
	}{ // Starting the test cases.
		{"Simple palindrome", "racecar", true},               // Case: basic palindrome.
		{"With uppercase", "RaceCar", true},                  // Case: case-insensitivity.
		{"With spaces", "A man a plan a canal Panama", true}, // Case: ignoring spaces.
		{"With punctuation", "No 'x' in Nixon", true},        // Case: ignoring punctuation.
		{"Not a palindrome", "hello", false},                 // Case: not a palindrome.
		{"Numbers", "12321", true},                           // Case: numeric palindrome.
		{"Empty", "", true},                                  // Case: empty string is a palindrome.
	} // End of test cases.

	for _, tt := range tests { // Iterating through the test table.
		t.Run(tt.name, func(t *testing.T) { // Running a subtest for each case.
			if result := IsPalindrome(tt.input); result != tt.expected { // Calling the function and checking result.
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, result, tt.expected) // Logging error on mismatch.
			} // End of result check.
		}) // End of subtest.
	} // End of test loop.
} // Closing the TestIsPalindrome function.

func TestSliceManipulation(t *testing.T) { // Test function for slice-specific behaviors.
	// Demonstrating slice properties through tests // Noting the purpose of these subtests.
	t.Run("Append and Capacity", func(t *testing.T) { // Subtest for append and capacity growth.
		s := make([]int, 0, 2)          // Creating a slice with length 0 and capacity 2.
		s = append(s, 1)                // Appending the first element.
		s = append(s, 2)                // Appending the second element.
		if len(s) != 2 || cap(s) != 2 { // Checking length and capacity.
			t.Errorf("Expected len 2, cap 2; got len %d, cap %d", len(s), cap(s)) // Logging error if incorrect.
		} // End of state check.

		s = append(s, 3)               // Appending a third element to trigger growth.
		if len(s) != 3 || cap(s) < 3 { // Checking that length grew and capacity increased.
			t.Errorf("Expected len 3, cap >= 3; got len %d, cap %d", len(s), cap(s)) // Logging error if incorrect.
		} // End of growth check.
	}) // End of subtest.

	t.Run("Slicing shares memory", func(t *testing.T) { // Subtest for memory sharing in slices.
		original := []int{1, 2, 3, 4, 5} // Initializing an original slice.
		sub := original[1:4]             // Creating a sub-slice {2, 3, 4}.

		sub[0] = 99            // Modifying the first element of the sub-slice.
		if original[1] != 99 { // Checking if the modification affected the original slice.
			t.Errorf("Expected original[1] to be 99 because sub-slice shares memory; got %d", original[1]) // Logging error on no effect.
		} // End of memory sharing check.
	}) // End of subtest.
} // Closing the TestSliceManipulation function.

func TestJSONMarshalling(t *testing.T) { // Test function for JSON processing.
	t.Run("Struct to JSON", func(t *testing.T) { // Subtest for struct-to-JSON conversion.
		p := Product{ID: 1, Name: "Test", Price: 10.5, Tags: []string{"a", "b"}} // Initializing a Product.
		data, err := json.Marshal(p)                                             // Marshalling the product to JSON.
		if err != nil {                                                          // Checking for marshalling errors.
			t.Fatalf("Marshal failed: %v", err) // Failing fast on error.
		} // End of error check.

		// Use a map to verify content without worrying about field order // Explaining the verification strategy.
		var m map[string]interface{}                     // Declaring a map to hold unmarshalled JSON.
		if err := json.Unmarshal(data, &m); err != nil { // Unmarshalling back to a map for verification.
			t.Fatalf("Unmarshal for verification failed: %v", err) // Failing fast on error.
		} // End of unmarshalling check.

		if m["id"] != 1.0 || m["name"] != "Test" { // Checking map values (JSON numbers are floats in Go).
			t.Errorf("JSON content mismatch: %v", m) // Logging error on mismatch.
		} // End of value check.
	}) // End of subtest.

	t.Run("Omitempty", func(t *testing.T) { // Subtest for the 'omitempty' JSON tag.
		p := Product{ID: 2, Name: "NoTags"} // Initializing a Product with empty Tags.
		data, err := json.Marshal(p)        // Marshalling the product.
		if err != nil {                     // Checking for marshalling errors.
			t.Fatalf("Marshal failed: %v", err) // Failing fast on error.
		} // End of error check.

		if strings.Contains(string(data), "tags") { // Checking that the "tags" key was omitted.
			t.Errorf("Expected 'tags' to be omitted, but found it in: %s", string(data)) // Logging error if found.
		} // End of omission check.
	}) // End of subtest.
} // Closing the TestJSONMarshalling function.

func TestMapOperations(t *testing.T) { // Test function for map-specific operations.
	t.Run("Comma ok idiom", func(t *testing.T) { // Subtest for checking map key existence.
		m := map[string]int{"a": 1} // Initializing a map.

		val, ok := m["a"]    // Checking for an existing key.
		if !ok || val != 1 { // Verifying existence and value.
			t.Errorf("Expected (1, true); got (%d, %t)", val, ok) // Logging error on incorrect result.
		} // End of existing key check.

		val, ok = m["b"]    // Checking for a non-existent key.
		if ok || val != 0 { // Verifying non-existence and zero-value.
			t.Errorf("Expected (0, false); got (%d, %t)", val, ok) // Logging error on incorrect result.
		} // End of missing key check.
	}) // End of subtest.

	t.Run("DeepEqual for maps", func(t *testing.T) { // Subtest for structural map equality.
		m1 := map[string]int{"a": 1, "b": 2} // First map.
		m2 := map[string]int{"b": 2, "a": 1} // Second map with different key order.

		if !reflect.DeepEqual(m1, m2) { // Checking if both maps are structurally equal.
			t.Error("Maps with same content should be DeepEqual") // Logging error on inequality.
		} // End of equality check.
	}) // End of subtest.
} // Closing the TestMapOperations function.
