// Package generics showcases Go's generic capabilities. // Package comment describing purpose.
package generics

import (
	"fmt"
)

// List is a generic slice wrapper. // Comment describing the List type.
//
// For a Java developer: // Guidance for Java background readers.
//   - Go generics use square brackets `[T]` instead of angle brackets `<T>`. // Generics syntax difference.
//   - `any` is a type constraint equivalent to `Object` in Java (but it's actually an alias for `interface{}`). // Clarifying `any`.
//   - NO Type Erasure: Go preserves type information at runtime for monomorphized types. // Highlighting no type erasure.
//     `List[int]` is a different type than `List[string]` at runtime. // Emphasizing runtime distinct types.
//   - Methods on generic types must also declare the type parameters. // Note about method declarations on generic types.
type List[T any] struct { // Declaring a generic struct List with type parameter T.
	items []T // Underlying slice storing items of type T.
}

// Add appends an item to the list. // Comment for the Add method.
// In Go, we use a pointer receiver `(l *List[T])` to modify the struct's internal state. // Explaining pointer receiver necessity.
func (l *List[T]) Add(item T) { // Method Add with pointer receiver to mutate List.
	l.items = append(l.items, item) // `append` is a built-in function that handles slice resizing and copying as needed.
}

// Get returns the item at the specified index and a boolean indicating success. // Comment for the Get method.
// Java equivalent: `public Optional<T> get(int index)` // Comparing to Java Optional.
// Go doesn't have `Optional`; it's idiomatic to return `(value, ok)`. // Explaining idiomatic pair return.
func (l *List[T]) Get(index int) (T, bool) { // Method Get returning value and presence flag.
	if index < 0 || index >= len(l.items) { // Bounds check to prevent panic on out-of-range.
		var zero T         // Returns the "zero value" for type T (e.g., 0 for int, "" for string, nil for pointers)
		return zero, false // Indicating failure to get a value at the index.
	}
	return l.items[index], true
}

// For a Java developer: // Guidance for Java background readers.
// - This is like `list.stream().map(f).collect(Collectors.toList())`. // Comparing to Java Stream map.
// - Go doesn't have a built-in Stream API, so we often write these helper functions. // Explaining helper usage.
func MapValues[T any, R any](items []T, f func(T) R) []R { // Generic MapValues transforming []T to []R.
	result := make([]R, len(items)) // Pre-allocating result with the same length as input for efficiency.
	for i, v := range items {       // Iterating over items with index and value.
		result[i] = f(v) // Applying the mapping function to each element and storing it at the same index.
	}
	return result
}

// For a Java developer: // Guidance for Java background readers.
// - Demonstrates using generic types and functions with different concrete types. // Explaining demo purpose.
// - Note how Go infers the type parameters in function calls like `MapValues(ints, ...)`. // Noting type inference in calls.
func RunGenericsDemo() { // Demo function printing examples to stdout.
	fmt.Println("--- Generics Demo ---") // Printing section header.

	// 1. Using a generic type // Section: generic type usage.
	// Java equivalent: List<Integer> intList = new ArrayList<>(); // Comparing to Java list instantiation.
	fmt.Println("1. Generic List type (int and string):") // Printing subsection title.
	intList := List[int]{}                                // Declaring a zero-value generic List of int.
	intList.Add(10)                                       // Adding first integer element.
	intList.Add(20)                                       // Adding second integer element.
	val, _ := intList.Get(0)                              // Getting the first element (ignoring ok for brevity).
	fmt.Printf("   Int List value: %d\n", val)            // Printing retrieved int value.

	// Java equivalent: List<String> stringList = new ArrayList<>(); // Comparing to Java string list.
	stringList := List[string]{}                   // Declaring a zero-value generic List of string.
	stringList.Add("Go")                           // Adding first string element.
	stringList.Add("Generics")                     // Adding second string element.
	sVal, _ := stringList.Get(1)                   // Getting the second element (index 1).
	fmt.Printf("   String List value: %s\n", sVal) // Printing retrieved string value.

	// 2. Using a generic function // Section: generic function usage.
	// Java equivalent: Stream.of(1, 2, 3).map(i -> i * 2).toList(); // Comparing to Java stream pipeline.
	fmt.Println("2. Generic function MapValues:") // Printing subsection title.
	ints := []int{1, 2, 3, 4, 5}                  // Defining a slice of integers.
	doubled := MapValues(ints, func(i int) int {  // Mapping each int to double its value.
		return i * 2
	})
	fmt.Printf("   Doubled: %v\n", doubled) // Printing the doubled slice.

	lengths := MapValues([]string{"Go", "Generics", "Showcase"}, func(s string) int { // Mapping strings to their lengths.
		return len(s)
	})
	fmt.Printf("   Word lengths: %v\n", lengths) // Printing resulting lengths.

	fmt.Println("--- Generics Demo End ---") // Printing section footer.
}
