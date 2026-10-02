// Package generics showcases Go's generic capabilities.
package generics

import (
	"fmt"
)

// List is a generic slice wrapper.
//
// For a Java developer:
//   - Go generics use square brackets `[T]` instead of angle brackets `<T>`.
//   - `any` is a type constraint equivalent to `Object` in Java (an alias for `interface{}`).
//   - NO Type Erasure: Go preserves type information at runtime for monomorphized types.
//     `List[int]` is a distinct type from `List[string]` at runtime.
//   - Methods on generic types must declare the type parameter: `(l *List[T])`.
//   - Pointer receivers `(l *List[T])` are used when methods modify internal state.
type List[T any] struct {
	items []T
}

// Add appends an item to the list.
// `append` is a built-in function that handles slice resizing and copying as needed.
func (l *List[T]) Add(item T) {
	l.items = append(l.items, item)
}

// Get returns the item at the specified index and a boolean indicating success.
//
// For a Java developer:
//   - Java equivalent: `public Optional<T> get(int index)`.
//   - Go doesn't have `Optional`; it is idiomatic to return `(value, ok)`.
//   - On failure, returns the type's zero value (e.g. 0 for int, "" for string, nil for pointers) and false.
func (l *List[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, false
	}
	return l.items[index], true
}

// MapValues applies a function to each element of a slice and returns a new slice.
//
// For a Java developer:
//   - Java equivalent: `list.stream().map(f).collect(Collectors.toList())`.
//   - Go lacks a built-in Stream API; generic helper functions are common for functional slice operations.
//   - Result slice is pre-allocated with `make([]R, len(items))` for efficiency.
func MapValues[T any, R any](items []T, f func(T) R) []R {
	result := make([]R, len(items))
	for i, v := range items {
		result[i] = f(v)
	}
	return result
}

// RunGenericsDemo showcases generic types and functions with different concrete types.
// Notice how Go infers type parameters in function calls like `MapValues(ints, ...)`.
func RunGenericsDemo() {
	fmt.Println("--- Generics Demo ---")

	// 1. Using a generic type (int and string)
	// Java equivalent: List<Integer> intList = new ArrayList<>();
	fmt.Println("1. Generic List type (int and string):")
	intList := List[int]{}
	intList.Add(10)
	intList.Add(20)
	val, _ := intList.Get(0)
	fmt.Printf("   Int List value: %d\n", val)

	// Java equivalent: List<String> stringList = new ArrayList<>();
	stringList := List[string]{}
	stringList.Add("Go")
	stringList.Add("Generics")
	sVal, _ := stringList.Get(1)
	fmt.Printf("   String List value: %s\n", sVal)

	// 2. Using a generic function (MapValues)
	// Java equivalent: Stream.of(1, 2, 3).map(i -> i * 2).toList();
	fmt.Println("2. Generic function MapValues:")
	ints := []int{1, 2, 3, 4, 5}
	doubled := MapValues(ints, func(i int) int {
		return i * 2
	})
	fmt.Printf("   Doubled: %v\n", doubled)

	lengths := MapValues([]string{"Go", "Generics", "Showcase"}, func(s string) int {
		return len(s)
	})
	fmt.Printf("   Word lengths: %v\n", lengths)

	fmt.Println("--- Generics Demo End ---")
}
