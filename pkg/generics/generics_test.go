package generics

import (
	"reflect"
	"testing"
)

/**
 * ===========================================================================
 * Generics Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Go generics allow you to write tests once for many types.
 * - `reflect.DeepEqual` is often used to compare slices/maps/structs (like `assertEquals` for objects).
 */

func TestMapValues(t *testing.T) {
	// 1. Using a generic function with different types
	t.Run("Ints to Ints", func(t *testing.T) {
		ints := []int{1, 2, 3}
		doubled := MapValues(ints, func(i int) int { return i * 2 })
		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(doubled, expected) {
			t.Errorf("Expected %v, got %v", expected, doubled)
		}
	})

	t.Run("Strings to Ints", func(t *testing.T) {
		strs := []string{"a", "ab", "abc"}
		lengths := MapValues(strs, func(s string) int { return len(s) })
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(lengths, expected) {
			t.Errorf("Expected %v, got %v", expected, lengths)
		}
	})
}

func TestGenericList(t *testing.T) {
	l := List[string]{}
	l.Add("first")
	l.Add("second")

	val, ok := l.Get(0)
	if !ok || val != "first" {
		t.Errorf("Expected 'first', got %v (ok: %v)", val, ok)
	}

	val, ok = l.Get(1)
	if !ok || val != "second" {
		t.Errorf("Expected 'second', got %v (ok: %v)", val, ok)
	}

	_, ok = l.Get(2)
	if ok {
		t.Errorf("Expected not ok for index 2")
	}
}
