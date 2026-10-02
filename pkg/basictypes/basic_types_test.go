package basictypes

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Basic Types & Standard Library Testing (For a Java developer)
//
// In Go, testing slices and maps often uses `reflect.DeepEqual` for structural
// comparison, as the `==` operator is not defined for slices and only for
// simple maps (comparing to nil).
//
// Table-driven tests are used here to test string manipulation and JSON logic.
// ---------------------------------------------------------------------------

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single character", "a", "a"},
		{"Simple word", "hello", "olleh"},
		{"Palindrome", "racecar", "racecar"},
		{"Multi-byte characters", "こんにちは", "はちにんこ"},
		{"Mixed", "Go 1.23", "32.1 oG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseString(tt.input)
			if result != tt.expected {
				t.Errorf("ReverseString(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Simple palindrome", "racecar", true},
		{"With uppercase", "RaceCar", true},
		{"With spaces", "A man a plan a canal Panama", true},
		{"With punctuation", "No 'x' in Nixon", true},
		{"Not a palindrome", "hello", false},
		{"Numbers", "12321", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := IsPalindrome(tt.input); result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSliceManipulation(t *testing.T) {
	// Subtest 1: Slices grow their capacity dynamically when appended
	t.Run("Append and Capacity", func(t *testing.T) {
		s := make([]int, 0, 2)
		s = append(s, 1)
		s = append(s, 2)
		if len(s) != 2 || cap(s) != 2 {
			t.Errorf("Expected len 2, cap 2; got len %d, cap %d", len(s), cap(s))
		}

		s = append(s, 3)
		if len(s) != 3 || cap(s) < 3 {
			t.Errorf("Expected len 3, cap >= 3; got len %d, cap %d", len(s), cap(s))
		}
	})

	// Subtest 2: Slices are views into the underlying array, so modifying a sub-slice affects the original
	t.Run("Slicing shares memory", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		sub := original[1:4]

		sub[0] = 99
		if original[1] != 99 {
			t.Errorf("Expected original[1] to be 99 because sub-slice shares memory; got %d", original[1])
		}
	})
}

func TestJSONMarshalling(t *testing.T) {
	// Verify struct to JSON serialization and content equality via unmarshalled map
	t.Run("Struct to JSON", func(t *testing.T) {
		p := Product{ID: 1, Name: "Test", Price: 10.5, Tags: []string{"a", "b"}}
		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("Unmarshal for verification failed: %v", err)
		}

		if m["id"] != 1.0 || m["name"] != "Test" {
			t.Errorf("JSON content mismatch: %v", m)
		}
	})

	// Verify omitempty tag omits zero-valued slice field
	t.Run("Omitempty", func(t *testing.T) {
		p := Product{ID: 2, Name: "NoTags"}
		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		if strings.Contains(string(data), "tags") {
			t.Errorf("Expected 'tags' to be omitted, but found it in: %s", string(data))
		}
	})
}

func TestMapOperations(t *testing.T) {
	// Verify the comma-ok idiom for key lookup
	t.Run("Comma ok idiom", func(t *testing.T) {
		m := map[string]int{"a": 1}

		val, ok := m["a"]
		if !ok || val != 1 {
			t.Errorf("Expected (1, true); got (%d, %t)", val, ok)
		}

		val, ok = m["b"]
		if ok || val != 0 {
			t.Errorf("Expected (0, false); got (%d, %t)", val, ok)
		}
	})

	// Maps are unordered, but reflect.DeepEqual correctly considers key-value equality
	t.Run("DeepEqual for maps", func(t *testing.T) {
		m1 := map[string]int{"a": 1, "b": 2}
		m2 := map[string]int{"b": 2, "a": 1}

		if !reflect.DeepEqual(m1, m2) {
			t.Error("Maps with same content should be DeepEqual")
		}
	})
}
