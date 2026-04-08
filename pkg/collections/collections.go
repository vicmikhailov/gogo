// Package collections provides generic data structures and functional slice operations
// comparable to Java Collections and Stream API.
//
// For a Java developer:
//   - Go prefers built-in types (slices and maps) over complex collection hierarchies.
//   - Java-ism to avoid: Creating custom collection types for everything.
//     In Go, a simple slice `[]T` is usually all you need.
//   - This package demonstrates how to build such structures when necessary,
//     but lean towards standard slices/maps in your own code.
package collections

import (
	"fmt"
	"gogo/pkg/generics"
	"sort"
	"strings"
)

// For a Java developer: // Explanation targeted at Java developers.
// - Go doesn't have a built-in `Set`. We usually use a `map[T]struct{}`. // Explaining the Go idiomatic Set implementation.
// - The `comparable` constraint ensures types can be used as map keys. // Noting the type constraint for keys.
type Set[T comparable] struct { // Defining a generic Set struct with comparable type T.
	items map[T]struct{} // Underlying map where keys are set elements and values are empty structs.
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `new HashSet<>(Arrays.asList(items))` // Comparing to Java HashSet creation.
// - Go variadic functions (`...T`) are like Java's `T...`. // Explaining Go's variadic parameter syntax.
func NewSet[T comparable](items ...T) *Set[T] { // NewSet function taking variadic arguments.
	s := &Set[T]{items: make(map[T]struct{})} // Initializing a new Set with an empty map.
	for _, item := range items {              // Iterating through the provided items.
		s.Add(item) // Adding each item to the set.
	}
	return s
}

// Add adds an item to the set. // Comment for the Add method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `set.add(item)` // Comparing to Java Set.add.
// - In Go, maps are always passed by reference (technically they are pointers). // Explaining map passing semantics.
func (s *Set[T]) Add(item T) { s.items[item] = struct{}{} } // Method Add with Set pointer receiver, inserting item into map.

// Remove removes an item from the set. // Comment for the Remove method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `set.remove(item)` // Comparing to Java Set.remove.
func (s *Set[T]) Remove(item T) { delete(s.items, item) } // Method Remove with Set pointer receiver, deleting item from map.

// Contains checks if an item is in the set. // Comment for the Contains method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `set.contains(item)` // Comparing to Java Set.contains.
// - The "comma ok" idiom (`_, ok := m[key]`) is the standard way to check for key existence in Go. // Explaining existence checks.
func (s *Set[T]) Contains(item T) bool { _, ok := s.items[item]; return ok } // Method Contains returning a boolean existence status.

// Len returns the number of items in the set. // Comment for the Len method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `set.size()` // Comparing to Java Set.size.
func (s *Set[T]) Len() int { return len(s.items) } // Method Len returning the number of entries in the map.

// Values returns all elements in arbitrary order as a slice. // Comment for the Values method.
// Java equivalent: `new ArrayList<>(set)` // Comparing to Java list creation from set.
func (s *Set[T]) Values() []T { // Method Values returning a slice of type T.
	result := make([]T, 0, len(s.items)) // Initializing a slice with 0 length and capacity equal to set size.
	for k := range s.items {             // Iterating through map keys.
		result = append(result, k) // Appending each key to the result slice.
	}
	return result
}

// Union returns a new set containing elements from both sets. // Comment for the Union method.
// Java equivalent: `Set<T> union = new HashSet<>(s); union.addAll(other);` // Comparing to Java set union.
func (s *Set[T]) Union(other *Set[T]) *Set[T] { // Method Union returning a new Set.
	result := NewSet[T]()    // Creating a new empty Set.
	for k := range s.items { // Iterating through elements of the first set.
		result.Add(k) // Adding element to the result.
	}
	for k := range other.items { // Iterating through elements of the second set.
		result.Add(k) // Adding element to the result (uniqueness handled by map).
	}
	return result
}

// Intersection returns a new set containing only elements present in both sets. // Comment for the Intersection method.
// Java equivalent: `Set<T> intersect = new HashSet<>(s); intersect.retainAll(other);` // Comparing to Java set intersection.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] { // Method Intersection returning a new Set.
	result := NewSet[T]()    // Creating a new empty Set.
	for k := range s.items { // Iterating through elements of the first set.
		if other.Contains(k) { // Checking if element exists in the second set.
			result.Add(k) // Adding to result if present in both.
		}
	}
	return result
}

// Difference returns a new set containing elements in s but not in other. // Comment for the Difference method.
// Java equivalent: `Set<T> diff = new HashSet<>(s); diff.removeAll(other);` // Comparing to Java set difference.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] { // Method Difference returning a new Set.
	result := NewSet[T]()    // Creating a new empty Set.
	for k := range s.items { // Iterating through elements of the first set.
		if !other.Contains(k) { // Checking if element is absent in the second set.
			result.Add(k) // Adding to result if only in the first set.
		}
	}
	return result
}

// For a Java developer: // Explanation targeted at Java developers.
// - Go uses slices (`[]T`) for dynamic arrays (similar to `ArrayList`). // Comparing Go slices to Java ArrayList.
// - Slices are lightweight descriptors pointing to an underlying array. // Explaining the nature of slices.
type Stack[T any] struct { // Defining a generic Stack struct with type T.
	items []T // Underlying slice to store stack elements.
}

// Push adds an item to the top of the stack. // Comment for the Push method.
// Java equivalent: `stack.push(item)` // Comparing to Java stack.push.
func (s *Stack[T]) Push(item T) { s.items = append(s.items, item) } // Method Push appending item to the end of the slice.

// Pop removes and returns the top item from the stack. // Comment for the Pop method.
// Returns the zero value and false if the stack is empty. // Specifying the return behavior for empty stack.
// Java equivalent: `stack.pop()` (but it throws an exception if empty) // Comparing to Java stack.pop.
func (s *Stack[T]) Pop() (T, bool) { // Method Pop returning value and success status.
	if len(s.items) == 0 { // Checking if the stack is empty.
		var zero T // Declaring a zero value of type T.
		return zero, false
	}
	item := s.items[len(s.items)-1]    // Retrieving the last element of the slice.
	s.items = s.items[:len(s.items)-1] // Reducing the slice length by one (O(1) operation).
	return item, true
}

// Peek returns the top item from the stack without removing it. // Comment for the Peek method.
// Returns the zero value and false if the stack is empty. // Specifying empty stack behavior.
// Java equivalent: `stack.peek()` // Comparing to Java stack.peek.
func (s *Stack[T]) Peek() (T, bool) { // Method Peek returning value and success status.
	if len(s.items) == 0 { // Checking if the stack is empty.
		var zero T // Declaring a zero value of type T.
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of items in the stack. // Comment for the Len method.
// Java equivalent: `stack.size()` // Comparing to Java stack.size.
func (s *Stack[T]) Len() int { return len(s.items) } // Method Len returning slice length.

// IsEmpty returns true if the stack contains no items. // Comment for the IsEmpty method.
// Java equivalent: `stack.isEmpty()` // Comparing to Java stack.isEmpty.
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 } // Method IsEmpty checking if slice length is zero.

// For a Java developer: // Explanation targeted at Java developers.
// - Comparable to Java LinkedList or ArrayDeque. // Comparing to Java queue implementations.
type Queue[T any] struct { // Defining a generic Queue struct with type T.
	items []T // Underlying slice to store queue elements.
}

// Enqueue adds an item to the end of the queue. // Comment for the Enqueue method.
// Java equivalent: `queue.offer(item)` or `queue.add(item)` // Comparing to Java queue insertion.
func (q *Queue[T]) Enqueue(item T) { q.items = append(q.items, item) } // Method Enqueue appending item to the end of the slice.

// Dequeue removes and returns the front item from the queue. // Comment for the Dequeue method.
// Returns the zero value and false if the queue is empty. // Specifying empty queue behavior.
// Java equivalent: `queue.poll()` // Comparing to Java queue.poll.
func (q *Queue[T]) Dequeue() (T, bool) { // Method Dequeue returning value and success status.
	if len(q.items) == 0 { // Checking if the queue is empty.
		var zero T // Declaring a zero value of type T.
		return zero, false
	}
	item := q.items[0]    // Retrieving the first element of the slice.
	q.items = q.items[1:] // Re-slicing to remove the first element (moves the internal pointer).
	return item, true
}

// Peek returns the front item from the queue without removing it. // Comment for the Peek method.
// Returns the zero value and false if the queue is empty. // Specifying empty queue behavior.
// Java equivalent: `queue.peek()` // Comparing to Java queue.peek.
func (q *Queue[T]) Peek() (T, bool) { // Method Peek returning value and success status.
	if len(q.items) == 0 { // Checking if the queue is empty.
		var zero T // Declaring a zero value of type T.
		return zero, false
	}
	return q.items[0], true
}

// Len returns the number of items in the queue. // Comment for the Len method.
// Java equivalent: `queue.size()` // Comparing to Java queue.size.
func (q *Queue[T]) Len() int { return len(q.items) } // Method Len returning slice length.

// IsEmpty returns true if the queue contains no items. // Comment for the IsEmpty method.
// Java equivalent: `queue.isEmpty()` // Comparing to Java queue.isEmpty.
func (q *Queue[T]) IsEmpty() bool { return len(q.items) == 0 } // Method IsEmpty checking if slice length is zero.

// For a Java developer: // Explanation targeted at Java developers.
// - Comparable to Java LinkedHashMap. // Comparing to Java LinkedHashMap.
type OrderedMap[K comparable, V any] struct {
	keys   []K     // Slice to maintain the order of keys.
	values map[K]V // Map to store key-value associations.
}

// For a Java developer: // Explanation targeted at Java developers.
//   - This is like calling `new LinkedHashMap<>()`. // Comparing to Java constructor call.
//   - Note that in Go, we often use factory functions (usually starting with `New`) // Explaining the Go factory pattern.
//     to initialize complex structs that require map/slice allocation. // Detail on factory function usage.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] { // NewOrderedMap function returning a pointer.
	return &OrderedMap[K, V]{values: make(map[K]V)} // Initializing and returning the struct with a new map.
}

// Put adds or updates a key-value pair. New keys are appended; existing keys keep their position. // Comment for the Put method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `map.put(key, value)` // Comparing to Java map.put.
// - Go's `map` doesn't track insertion order. We maintain a separate `keys` slice for this. // Explaining how order is maintained.
func (m *OrderedMap[K, V]) Put(key K, value V) { // Method Put with OrderedMap pointer receiver.
	if _, exists := m.values[key]; !exists { // Checking if the key already exists in the values map.
		m.keys = append(m.keys, key) // Appending the new key to the keys slice if it's new.
	}
	m.values[key] = value // Storing the value in the map.
}

// Get retrieves a value by key. // Comment for the Get method.
//
// For a Java developer: // Explanation targeted at Java developers.
//   - Java equivalent: `map.get(key)` // Comparing to Java map.get.
//   - Returns the value and a boolean `ok` indicating if the key was found. // Explaining the return behavior.
//   - This is safer than returning `null` because it distinguishes between "not found" // Explaining safety over Java's null.
//     and "value is actually the zero-value". // Clarifying why boolean ok is useful.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) { // Method Get returning value and existence status.
	v, ok := m.values[key] // Retrieving the value from the underlying map.
	return v, ok
}

// Delete removes a key-value pair. // Comment for the Delete method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `map.remove(key)` // Comparing to Java map.remove.
// - This is an O(N) operation because we have to find and remove the key from the `keys` slice. // Noting the time complexity.
func (m *OrderedMap[K, V]) Delete(key K) { // Method Delete with OrderedMap pointer receiver.
	if _, exists := m.values[key]; !exists { // Checking if the key exists before proceeding.
		return // Exiting early if key doesn't exist.
	}
	delete(m.values, key)      // Deleting the key from the underlying map.
	for i, k := range m.keys { // Iterating through the keys slice to find the key to remove.
		if k == key { // Checking for a key match.
			// Efficient slice element removal: [0:i] + [i+1:end] // Comment explaining the slice removal technique.
			// The `...` operator is like the spread operator in JavaScript, // Explaining the variadic operator.
			// or passing elements individually from a collection in Java. // Comparing to Java equivalent.
			m.keys = append(m.keys[:i], m.keys[i+1:]...) // Removing the element by joining slices before and after it.
			break                                        // Exiting the loop once the key is removed.
		}
	}
}

// Keys returns keys in insertion order. // Comment for the Keys method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `new ArrayList<>(map.keySet())` // Comparing to Java keySet conversion.
// - We return a copy of the keys slice to prevent external modification of internal state. // Explaining why a copy is returned.
func (m *OrderedMap[K, V]) Keys() []K { // Method Keys returning a slice of keys.
	result := make([]K, len(m.keys)) // Creating a new slice to hold the copy.
	copy(result, m.keys)             // Copying internal keys to the result slice.
	return result
}

// Len returns the number of key-value pairs. // Comment for the Len method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `map.size()` // Comparing to Java map.size.
func (m *OrderedMap[K, V]) Len() int { return len(m.keys) } // Method Len returning the count of keys.

// ForEach iterates over entries in insertion order. // Comment for the ForEach method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `map.forEach((k, v) -> ...)` // Comparing to Java map.forEach.
// - In Go, passing a function as an argument is a common way to implement internal iteration. // Explaining the internal iteration pattern.
func (m *OrderedMap[K, V]) ForEach(fn func(K, V)) {
	for _, k := range m.keys { // Iterating through the keys in insertion order.
		fn(k, m.values[k]) // Calling the callback function with each key and its value.
	}
}

// For a Java developer: // Explanation targeted at Java developers.
// - Go doesn't have a built-in `Stream` type. // Noting the lack of a Stream API.
// - This is like `list.stream().filter(predicate).collect(Collectors.toList())`. // Comparing to Java Stream filter.
// - The `make([]T, 0)` initializes an empty slice. // Noting how empty slices are initialized.
func Filter[T any](items []T, predicate func(T) bool) []T { // Generic Filter function taking a slice and a predicate.
	result := make([]T, 0)    // Initializing an empty result slice.
	for _, v := range items { // Iterating through each item in the input slice.
		if predicate(v) { // Applying the predicate to the item.
			result = append(result, v) // Appending the item to result if it matches the predicate.
		}
	}
	return result
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().reduce(initial, accumulator)` // Comparing to Java Stream reduce.
// - Go doesn't have a method-based Stream API, so we use standalone generic functions. // Explaining the use of standalone functions.
func Reduce[T any, R any](items []T, initial R, fn func(R, T) R) R { // Generic Reduce function with initial value and accumulator.
	acc := initial            // Initializing the accumulator with the provided initial value.
	for _, v := range items { // Iterating through each item in the slice.
		acc = fn(acc, v)
	}
	return acc
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().flatMap(fn).toList()` // Comparing to Java Stream flatMap.
// - The `...` operator in `append(result, fn(v)...)` expands the slice returned by `fn(v)`. // Explaining the variadic expansion operator.
func FlatMap[T any, R any](items []T, fn func(T) []R) []R {
	result := make([]R, 0)    // Initializing an empty result slice of type R.
	for _, v := range items { // Iterating through each item in the input slice.
		result = append(result, fn(v)...) // Applying the mapping function and flattening the resulting slice into result.
	}
	return result
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().collect(Collectors.groupingBy(keyFn))` // Comparing to Java groupingBy.
// - In Go, we use a `map[K][]T` where `K` must be `comparable`. // Explaining the Go group storage structure.
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T { // Generic GroupBy function returning a map.
	result := make(map[K][]T) // Initializing a map where values are slices of items.
	for _, v := range items { // Iterating through each item in the input slice.
		key := keyFn(v)                      // Computing the group key for the item.
		result[key] = append(result[key], v) // Adding the item to the appropriate group slice in the map.
	}
	return result
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().collect(Collectors.partitioningBy(predicate))` // Comparing to Java partitioningBy.
// - Go's ability to return multiple values makes this very clean. // Highlighting Go's multiple return values feature.
func Partition[T any](items []T, predicate func(T) bool) (matching []T, rest []T) { // Generic Partition function returning two slices.
	for _, v := range items { // Iterating through each item in the input slice.
		if predicate(v) { // Checking if the item matches the predicate.
			matching = append(matching, v) // Adding to matching slice if true.
		} else { // If the item doesn't match the predicate.
			rest = append(rest, v) // Adding to rest slice if false.
		}
	}
	return
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().sorted(comparator).toList()` // Comparing to Java Stream sorted.
// - Note: `sort.Slice` sorts in-place, so we copy the slice first. // Warning about in-place sorting and explaining the copy.
func Sorted[T any](items []T, less func(a, b T) bool) []T { // Generic Sorted function with comparator.
	result := make([]T, len(items))                                               // Creating a new slice to hold the sorted copy.
	copy(result, items)                                                           // Copying original items to the new slice.
	sort.Slice(result, func(i, j int) bool { return less(result[i], result[j]) }) // Sorting the copy in-place.
	return result
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().distinct().toList()` // Comparing to Java Stream distinct.
// - We use a `map[T]struct{}` as a temporary set to track seen elements. // Explaining how uniqueness is tracked.
func Distinct[T comparable](items []T) []T { // Generic Distinct function for comparable types.
	seen := make(map[T]struct{}) // Initializing a map to track encountered elements.
	result := make([]T, 0)       // Initializing an empty result slice.
	for _, v := range items {    // Iterating through each item in the input slice.
		if _, ok := seen[v]; !ok { // Checking if the item has not been seen yet.
			seen[v] = struct{}{}       // Marking the item as seen.
			result = append(result, v) // Adding the item to the result slice.
		}
	}
	return result
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().anyMatch(predicate)` // Comparing to Java Stream anyMatch.
func Any[T any](items []T, predicate func(T) bool) bool { // Generic Any function returning a boolean.
	for _, v := range items { // Iterating through each item in the input slice.
		if predicate(v) { // Checking if the current item matches the predicate.
			return true
		}
	}
	return false
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `list.stream().allMatch(predicate)` // Comparing to Java Stream allMatch.
func All[T any](items []T, predicate func(T) bool) bool { // Generic All function returning a boolean.
	for _, v := range items { // Iterating through each item in the input slice.
		if !predicate(v) { // Checking if any item fails the predicate.
			return false
		}
	}
	return true
}

// For a Java developer: // Explanation targeted at Java developers.
// - This is like `Map.Entry<A, B>` or a custom `Pair<A, B>` class. // Comparing to Java entry or pair classes.
// - Go doesn't have a built-in `Pair` or `Tuple` (besides multi-value returns). // Noting the lack of built-in pairs.
type Pair[A any, B any] struct { // Defining a generic Pair struct with types A and B.
	First  A // The first value of the pair.
	Second B // The second value of the pair.
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java doesn't have a built-in `zip` in the Stream API (often provided by Guava or Vavr). // Noting the lack of built-in zip in Java.
// - This demonstrates Go's generic functions handling multiple type parameters. // Highlighting multi-parameter generics.
func Zip[A any, B any](as []A, bs []B) []Pair[A, B] { // Generic Zip function taking two slices.
	minLen := len(as)     // Starting with the length of the first slice as minimum.
	if len(bs) < minLen { // Checking if the second slice is shorter.
		minLen = len(bs) // Updating minimum length if necessary.
	}
	result := make([]Pair[A, B], minLen) // Initializing the result slice with the determined minimum length.
	for i := 0; i < minLen; i++ {        // Iterating up to the minimum length.
		result[i] = Pair[A, B]{First: as[i], Second: bs[i]} // Creating and storing a Pair for each index.
	}
	return result
}

// RunCollectionsDemo showcases generic data structures and functional operations // Comment for the demo runner.
// comparable to Java Collections Framework and Stream API. // Noting the scope of the demo.
func RunCollectionsDemo() {
	fmt.Println("--- Collections & Data Structures Demo ---") // Printing the overall section header.
	intLess := func(x, y int) bool { return x < y }           // Defining a local comparator function for integers.

	// 1. Set operations // Comment for set operations demo.
	// Java equivalent: HashSet<Integer> a = new HashSet<>(Arrays.asList(1, 2, 3, 4, 5)); // Comparing to Java set initialization.
	fmt.Println("1. Generic Set (like Java HashSet):")                               // Printing section title.
	a := NewSet(1, 2, 3, 4, 5)                                                       // Initializing set A.
	b := NewSet(4, 5, 6, 7, 8)                                                       // Initializing set B.
	fmt.Printf("   Set A:        %v\n", Sorted(a.Values(), intLess))                 // Printing sorted values of A.
	fmt.Printf("   Set B:        %v\n", Sorted(b.Values(), intLess))                 // Printing sorted values of B.
	fmt.Printf("   Union:        %v\n", Sorted(a.Union(b).Values(), intLess))        // Printing sorted union of A and B.
	fmt.Printf("   Intersection: %v\n", Sorted(a.Intersection(b).Values(), intLess)) // Printing sorted intersection of A and B.
	fmt.Printf("   Difference:   %v\n", Sorted(a.Difference(b).Values(), intLess))   // Printing sorted difference (A - B).

	// 2. Stack (LIFO) // Comment for stack demo.
	// Java equivalent: Deque<String> stack = new ArrayDeque<>(); // Comparing to Java stack implementations.
	fmt.Println("2. Generic Stack (LIFO):") // Printing section title.
	stack := Stack[string]{}                // Initializing an empty generic stack of strings.
	stack.Push("first")                     // Pushing "first" onto stack.
	stack.Push("second")                    // Pushing "second".
	stack.Push("third")                     // Pushing "third".
	for !stack.IsEmpty() {                  // Iterating until the stack is empty.
		val, _ := stack.Pop()              // Popping the top element.
		fmt.Printf("   Popped: %s\n", val) // Printing the popped element.
	}

	// 3. Queue (FIFO) // Comment for queue demo.
	// Java equivalent: Queue<String> queue = new LinkedList<>(); // Comparing to Java queue implementations.
	fmt.Println("3. Generic Queue (FIFO):") // Printing section title.
	queue := Queue[string]{}                // Initializing an empty generic queue of strings.
	queue.Enqueue("task-1")                 // Enqueueing task-1.
	queue.Enqueue("task-2")                 // Enqueueing task-2.
	queue.Enqueue("task-3")                 // Enqueueing task-3.
	for !queue.IsEmpty() {                  // Iterating until the queue is empty.
		val, _ := queue.Dequeue()            // Dequeueing the front element.
		fmt.Printf("   Dequeued: %s\n", val) // Printing the dequeued element.
	}

	// 4. OrderedMap (like Java LinkedHashMap) // Comment for ordered map demo.
	fmt.Println("4. OrderedMap (insertion-order preserving, like LinkedHashMap):") // Printing section title.
	om := NewOrderedMap[string, int]()                                             // Initializing a new OrderedMap.
	om.Put("charlie", 3)                                                           // Adding "charlie" -> 3.
	om.Put("alpha", 1)                                                             // Adding "alpha" -> 1.
	om.Put("bravo", 2)                                                             // Adding "bravo" -> 2.
	om.ForEach(func(k string, v int) {                                             // Iterating over entries in insertion order.
		fmt.Printf("   %s -> %d\n", k, v) // Printing each entry.
	})

	// 5. Functional slice operations (like Java Streams) // Comment for functional operations demo.
	// Java equivalent: numbers.stream().filter(n -> n % 2 == 0).toList(); // Comparing to Java Stream filtering.
	fmt.Println("5. Functional operations (Java Stream equivalent):") // Printing section title.
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}                   // Defining a slice of integers.
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })    // Filtering for even numbers.
	fmt.Printf("   Filter (evens): %v\n", evens)                      // Printing filtered result.

	sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n }) // Reducing numbers to their sum.
	fmt.Printf("   Reduce (sum):   %d\n", sum)                         // Printing the sum.

	doubled := generics.MapValues(numbers, func(n int) int { return n * 2 }) // Mapping numbers to their doubled values.
	fmt.Printf("   Map (doubled):  %v\n", doubled)                           // Printing the doubled values.

	words := []string{"hello world", "go lang"}                                        // Defining a slice of strings.
	tokens := FlatMap(words, func(s string) []string { return strings.Split(s, " ") }) // Flat-mapping strings to their words.
	fmt.Printf("   FlatMap (split): %v\n", tokens)                                     // Printing the flattened words.

	// 6. GroupBy (like Collectors.groupingBy) // Comment for GroupBy demo.
	fmt.Println("6. GroupBy (Collectors.groupingBy equivalent):") // Printing section title.
	type Person struct {                                          // Defining a local Person struct for the demo.
		Name string // Person's name.
		City string // Person's city.
	}
	people := []Person{ // Initializing a slice of people.
		{"Alice", "NYC"}, {"Bob", "LA"}, {"Charlie", "NYC"}, // Group 1.
		{"Diana", "LA"}, {"Eve", "Chicago"}, // Group 2.
	}
	byCity := GroupBy(people, func(p Person) string { return p.City }) // Grouping people by their city.
	// Print in deterministic order for readability // Explaining the need for sorting keys.
	cityKeys := make([]string, 0, len(byCity)) // Creating a slice to hold city keys.
	for city := range byCity {                 // Collecting keys from the group map.
		cityKeys = append(cityKeys, city) // Appending city name to keys slice.
	}
	sort.Strings(cityKeys)          // Sorting city names alphabetically.
	for _, city := range cityKeys { // Iterating through sorted cities.
		names := generics.MapValues(byCity[city], func(p Person) string { return p.Name }) // Extracting names for each person in the city.
		fmt.Printf("   %s: %v\n", city, names)                                             // Printing the city and the list of names.
	}

	// 7. Chaining operations (pipeline style, like Stream pipelines) // Comment for chained operations demo.
	// Java equivalent: numbers.stream().filter(e).map(s).reduce(r); // Comparing to Java Stream chaining.
	fmt.Println("7. Chained pipeline (filter -> map -> reduce):") // Printing section title.
	result := Reduce(
		generics.MapValues( // Intermediate mapping step.
			Filter(numbers, func(n int) bool { return n%2 == 0 }), // Initial filtering step.
			func(n int) int { return n * n },                      // Squaring each even number.
		),
		0,                                       // Initial accumulator value for reduction.
		func(acc, n int) int { return acc + n }, // Accumulating sum of squares.
	)
	fmt.Printf("   Sum of squares of evens (2²+4²+6²+8²+10²): %d\n", result) // Printing final pipeline result.

	// 8. Partition, Distinct, Any, All // Comment for miscellaneous operations demo.
	fmt.Println("8. Partition, Distinct, Any, All:")                                  // Printing section title.
	pos, neg := Partition([]int{-3, -1, 0, 2, 5}, func(n int) bool { return n >= 0 }) // Partitioning into positive and negative.
	fmt.Printf("   Partition (>=0): %v | (<0): %v\n", pos, neg)                       // Printing partitioned result.

	unique := Distinct([]int{1, 2, 2, 3, 3, 3, 4}) // Extracting unique elements.
	fmt.Printf("   Distinct: %v\n", unique)        // Printing unique elements.

	hasNeg := Any(numbers, func(n int) bool { return n < 0 }) // Checking if any number is negative.
	fmt.Printf("   Any negative in 1..10? %v\n", hasNeg)      // Printing the result.

	allPos := All(numbers, func(n int) bool { return n > 0 }) // Checking if all numbers are positive.
	fmt.Printf("   All positive in 1..10? %v\n", allPos)      // Printing the result.

	// 9. Zip // Comment for zip operation demo.
	fmt.Println("9. Zip (combine two slices into pairs):") // Printing section title.
	names := []string{"Alice", "Bob", "Charlie"}           // Defining a slice of names.
	ages := []int{30, 25, 35}                              // Defining a slice of ages.
	pairs := Zip(names, ages)                              // Zipping names and ages into pairs.
	for _, p := range pairs {                              // Iterating through the pairs.
		fmt.Printf("   %s is %d years old\n", p.First, p.Second) // Printing each person's age.
	}

	fmt.Println("--- Collections & Data Structures Demo End ---") // Printing section footer.
}
