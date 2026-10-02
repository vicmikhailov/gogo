// Package collections provides generic data structures and functional slice operations
// comparable to Java Collections and the Stream API.
//
// For a Java developer:
//   - Go prefers built-in types (slices and maps) over complex collection hierarchies.
//   - Java-ism to avoid: Creating custom collection types for everything.
//     In Go, a simple slice `[]T` or map `map[K]V` is usually all you need.
//   - This package demonstrates how to build such structures when necessary,
//     but lean towards standard slices/maps in your own code.
package collections

import (
	"fmt"
	"gogo/pkg/generics"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// 1. Generic Set
// ---------------------------------------------------------------------------

// Set represents an unordered collection of unique elements.
//
// For a Java developer:
//   - Go doesn't have a built-in `Set`. The idiomatic implementation is `map[T]struct{}`.
//   - An empty struct `struct{}` takes zero bytes of memory.
//   - The `comparable` constraint ensures elements can be used as map keys.
//   - In Go, maps are reference types (pointers to internal runtime structures).
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates and initializes a Set with optional initial items.
// Java equivalent: `new HashSet<>(Arrays.asList(items))`
// Go variadic arguments (`...T`) work like Java's `T...`.
func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{items: make(map[T]struct{})}
	for _, item := range items {
		s.Add(item)
	}
	return s
}

// Add inserts an item into the set.
// Java equivalent: `set.add(item)`
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove deletes an item from the set.
// Java equivalent: `set.remove(item)`
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains reports whether item is present in the set using the "comma ok" idiom.
// Java equivalent: `set.contains(item)`
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}

// Len returns the number of elements in the set.
// Java equivalent: `set.size()`
func (s *Set[T]) Len() int {
	return len(s.items)
}

// Values returns all set elements in arbitrary order as a slice.
// Java equivalent: `new ArrayList<>(set)`
func (s *Set[T]) Values() []T {
	result := make([]T, 0, len(s.items))
	for k := range s.items {
		result = append(result, k)
	}
	return result
}

// Union returns a new set containing all elements from both s and other.
// Java equivalent: `Set<T> union = new HashSet<>(s); union.addAll(other);`
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.items {
		result.Add(k)
	}
	for k := range other.items {
		result.Add(k)
	}
	return result
}

// Intersection returns a new set containing only elements present in both s and other.
// Java equivalent: `Set<T> intersect = new HashSet<>(s); intersect.retainAll(other);`
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.items {
		if other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// Difference returns a new set containing elements in s but not in other.
// Java equivalent: `Set<T> diff = new HashSet<>(s); diff.removeAll(other);`
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.items {
		if !other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// 2. Generic Stack (LIFO)
// ---------------------------------------------------------------------------

// Stack implements a generic LIFO (Last-In-First-Out) data structure backed by a slice.
//
// For a Java developer:
//   - Go uses slices (`[]T`) for dynamic arrays (similar to `ArrayList` or `ArrayDeque`).
//   - Slices are lightweight headers pointing to an underlying array.
//   - Zero-value usable: `Stack[T]{}` is ready to use without an explicit constructor.
type Stack[T any] struct {
	items []T
}

// Push adds an item to the top of the stack.
// Java equivalent: `stack.push(item)`
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack.
// Returns the zero-value of T and false if the stack is empty.
// Java equivalent: `stack.pop()` (which throws NoSuchElementException if empty).
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// Peek returns the top item without removing it.
// Returns the zero-value of T and false if the stack is empty.
// Java equivalent: `stack.peek()`
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of elements in the stack.
// Java equivalent: `stack.size()`
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty reports whether the stack contains no elements.
// Java equivalent: `stack.isEmpty()`
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// ---------------------------------------------------------------------------
// 3. Generic Queue (FIFO)
// ---------------------------------------------------------------------------

// Queue implements a generic FIFO (First-In-First-Out) data structure backed by a slice.
//
// For a Java developer:
//   - Comparable to Java `Queue<T>` implemented by `LinkedList` or `ArrayDeque`.
//   - Zero-value usable: `Queue[T]{}` is ready to use immediately.
type Queue[T any] struct {
	items []T
}

// Enqueue adds an item to the end of the queue.
// Java equivalent: `queue.offer(item)` or `queue.add(item)`
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue.
// Returns the zero-value of T and false if the queue is empty.
// Java equivalent: `queue.poll()`
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the front item without removing it.
// Returns the zero-value of T and false if the queue is empty.
// Java equivalent: `queue.peek()`
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Len returns the number of elements in the queue.
// Java equivalent: `queue.size()`
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// IsEmpty reports whether the queue contains no elements.
// Java equivalent: `queue.isEmpty()`
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// ---------------------------------------------------------------------------
// 4. OrderedMap (Insertion-order preserving map)
// ---------------------------------------------------------------------------

// OrderedMap maintains key insertion order alongside O(1) map lookups.
//
// For a Java developer:
//   - Comparable to Java `LinkedHashMap<K, V>`.
//   - Standard Go maps do not preserve insertion order (iteration order is deliberately randomized).
//   - We maintain a separate `keys` slice to preserve insertion order.
type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

// NewOrderedMap creates an initialized OrderedMap.
// Java equivalent: `new LinkedHashMap<>()`
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{values: make(map[K]V)}
}

// Put adds or updates a key-value pair. New keys are appended; existing keys keep their position.
// Java equivalent: `map.put(key, value)`
func (m *OrderedMap[K, V]) Put(key K, value V) {
	if _, exists := m.values[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

// Get retrieves a value by key.
// Returns the value and a boolean indicating whether the key was found.
// Java equivalent: `map.get(key)`. Returning `(V, bool)` is safer than returning null because
// it distinguishes between a missing key and a present key whose value is the zero value.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.values[key]
	return v, ok
}

// Delete removes a key-value pair in O(N) time.
// Java equivalent: `map.remove(key)`
func (m *OrderedMap[K, V]) Delete(key K) {
	if _, exists := m.values[key]; !exists {
		return
	}
	delete(m.values, key)
	for i, k := range m.keys {
		if k == key {
			// Slice removal trick: append elements before index i with elements after index i
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

// Keys returns a copy of keys in insertion order.
// Java equivalent: `new ArrayList<>(map.keySet())`
func (m *OrderedMap[K, V]) Keys() []K {
	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return result
}

// Len returns the number of key-value pairs.
// Java equivalent: `map.size()`
func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

// ForEach iterates over entries in insertion order, invoking fn for each pair.
// Java equivalent: `map.forEach((k, v) -> ...)`
func (m *OrderedMap[K, V]) ForEach(fn func(K, V)) {
	for _, k := range m.keys {
		fn(k, m.values[k])
	}
}

// ---------------------------------------------------------------------------
// 5. Functional Slice Operations (Stream API equivalents)
// ---------------------------------------------------------------------------

// Filter returns a new slice containing only elements that satisfy the predicate.
// Java equivalent: `list.stream().filter(predicate).collect(Collectors.toList())`
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range items {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce aggregates elements using an accumulator function and initial value.
// Java equivalent: `list.stream().reduce(initial, accumulator)`
func Reduce[T any, R any](items []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, v := range items {
		acc = fn(acc, v)
	}
	return acc
}

// FlatMap maps each element to a slice and flattens the result into a single slice.
// Java equivalent: `list.stream().flatMap(fn).toList()`
// The `...` operator in `append(result, fn(v)...)` unpacks the slice.
func FlatMap[T any, R any](items []T, fn func(T) []R) []R {
	result := make([]R, 0)
	for _, v := range items {
		result = append(result, fn(v)...)
	}
	return result
}

// GroupBy groups elements by a key function into a map of slices.
// Java equivalent: `list.stream().collect(Collectors.groupingBy(keyFn))`
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range items {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Partition splits a slice into two slices based on a predicate:
// matching elements (where predicate is true) and rest (where predicate is false).
// Java equivalent: `list.stream().collect(Collectors.partitioningBy(predicate))`
func Partition[T any](items []T, predicate func(T) bool) (matching []T, rest []T) {
	for _, v := range items {
		if predicate(v) {
			matching = append(matching, v)
		} else {
			rest = append(rest, v)
		}
	}
	return
}

// Sorted returns a new sorted slice using the provided comparison function `less`.
// Java equivalent: `list.stream().sorted(comparator).toList()`
// Note: `sort.Slice` operates in-place, so we copy the slice first to avoid mutating input.
func Sorted[T any](items []T, less func(a, b T) bool) []T {
	result := make([]T, len(items))
	copy(result, items)
	sort.Slice(result, func(i, j int) bool { return less(result[i], result[j]) })
	return result
}

// Distinct returns a new slice containing only unique elements, preserving first-seen order.
// Java equivalent: `list.stream().distinct().toList()`
// Uses `map[T]struct{}` as an efficient lookup set.
func Distinct[T comparable](items []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, v := range items {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Any reports whether at least one element satisfies the predicate (short-circuiting).
// Java equivalent: `list.stream().anyMatch(predicate)`
func Any[T any](items []T, predicate func(T) bool) bool {
	for _, v := range items {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All reports whether all elements satisfy the predicate (short-circuiting).
// Java equivalent: `list.stream().allMatch(predicate)`
func All[T any](items []T, predicate func(T) bool) bool {
	for _, v := range items {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Pair holds two values of arbitrary types.
// Comparable to Java's `Map.Entry<A, B>` or a custom `Pair<A, B>` record.
type Pair[A any, B any] struct {
	First  A
	Second B
}

// Zip combines two slices element-wise into pairs up to the shorter slice's length.
// Java equivalent: Stream zip operations (provided by Guava or Vavr).
func Zip[A any, B any](as []A, bs []B) []Pair[A, B] {
	minLen := len(as)
	if len(bs) < minLen {
		minLen = len(bs)
	}
	result := make([]Pair[A, B], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = Pair[A, B]{First: as[i], Second: bs[i]}
	}
	return result
}

// RunCollectionsDemo showcases generic data structures and functional slice operations
// comparable to Java Collections Framework and Stream API.
func RunCollectionsDemo() {
	fmt.Println("--- Collections & Data Structures Demo ---")
	intLess := func(x, y int) bool { return x < y }

	// 1. Set operations (Union, Intersection, Difference)
	// Java equivalent: HashSet<Integer> a = new HashSet<>(Arrays.asList(1, 2, 3, 4, 5));
	fmt.Println("1. Generic Set (like Java HashSet):")
	a := NewSet(1, 2, 3, 4, 5)
	b := NewSet(4, 5, 6, 7, 8)
	fmt.Printf("   Set A:        %v\n", Sorted(a.Values(), intLess))
	fmt.Printf("   Set B:        %v\n", Sorted(b.Values(), intLess))
	fmt.Printf("   Union:        %v\n", Sorted(a.Union(b).Values(), intLess))
	fmt.Printf("   Intersection: %v\n", Sorted(a.Intersection(b).Values(), intLess))
	fmt.Printf("   Difference:   %v\n", Sorted(a.Difference(b).Values(), intLess))

	// 2. Stack (LIFO)
	// Java equivalent: Deque<String> stack = new ArrayDeque<>();
	fmt.Println("2. Generic Stack (LIFO):")
	stack := Stack[string]{}
	stack.Push("first")
	stack.Push("second")
	stack.Push("third")
	for !stack.IsEmpty() {
		val, _ := stack.Pop()
		fmt.Printf("   Popped: %s\n", val)
	}

	// 3. Queue (FIFO)
	// Java equivalent: Queue<String> queue = new LinkedList<>();
	fmt.Println("3. Generic Queue (FIFO):")
	queue := Queue[string]{}
	queue.Enqueue("task-1")
	queue.Enqueue("task-2")
	queue.Enqueue("task-3")
	for !queue.IsEmpty() {
		val, _ := queue.Dequeue()
		fmt.Printf("   Dequeued: %s\n", val)
	}

	// 4. OrderedMap (like Java LinkedHashMap)
	fmt.Println("4. OrderedMap (insertion-order preserving, like LinkedHashMap):")
	om := NewOrderedMap[string, int]()
	om.Put("charlie", 3)
	om.Put("alpha", 1)
	om.Put("bravo", 2)
	om.ForEach(func(k string, v int) {
		fmt.Printf("   %s -> %d\n", k, v)
	})

	// 5. Functional slice operations (Filter, Reduce, Map, FlatMap)
	// Java equivalent: numbers.stream().filter(n -> n % 2 == 0).toList();
	fmt.Println("5. Functional operations (Java Stream equivalent):")
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("   Filter (evens): %v\n", evens)

	sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
	fmt.Printf("   Reduce (sum):   %d\n", sum)

	doubled := generics.MapValues(numbers, func(n int) int { return n * 2 })
	fmt.Printf("   Map (doubled):  %v\n", doubled)

	words := []string{"hello world", "go lang"}
	tokens := FlatMap(words, func(s string) []string { return strings.Split(s, " ") })
	fmt.Printf("   FlatMap (split): %v\n", tokens)

	// 6. GroupBy (like Collectors.groupingBy)
	fmt.Println("6. GroupBy (Collectors.groupingBy equivalent):")
	type Person struct {
		Name string
		City string
	}
	people := []Person{
		{"Alice", "NYC"}, {"Bob", "LA"}, {"Charlie", "NYC"},
		{"Diana", "LA"}, {"Eve", "Chicago"},
	}
	byCity := GroupBy(people, func(p Person) string { return p.City })
	// Print in deterministic order by sorting city keys
	cityKeys := make([]string, 0, len(byCity))
	for city := range byCity {
		cityKeys = append(cityKeys, city)
	}
	sort.Strings(cityKeys)
	for _, city := range cityKeys {
		names := generics.MapValues(byCity[city], func(p Person) string { return p.Name })
		fmt.Printf("   %s: %v\n", city, names)
	}

	// 7. Chaining operations (pipeline style, like Stream pipelines)
	// Java equivalent: numbers.stream().filter(e).map(s).reduce(r);
	fmt.Println("7. Chained pipeline (filter -> map -> reduce):")
	result := Reduce(
		generics.MapValues(
			Filter(numbers, func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * n },
		),
		0,
		func(acc, n int) int { return acc + n },
	)
	fmt.Printf("   Sum of squares of evens (2²+4²+6²+8²+10²): %d\n", result)

	// 8. Partition, Distinct, Any, All
	fmt.Println("8. Partition, Distinct, Any, All:")
	pos, neg := Partition([]int{-3, -1, 0, 2, 5}, func(n int) bool { return n >= 0 })
	fmt.Printf("   Partition (>=0): %v | (<0): %v\n", pos, neg)

	unique := Distinct([]int{1, 2, 2, 3, 3, 3, 4})
	fmt.Printf("   Distinct: %v\n", unique)

	hasNeg := Any(numbers, func(n int) bool { return n < 0 })
	fmt.Printf("   Any negative in 1..10? %v\n", hasNeg)

	allPos := All(numbers, func(n int) bool { return n > 0 })
	fmt.Printf("   All positive in 1..10? %v\n", allPos)

	// 9. Zip
	fmt.Println("9. Zip (combine two slices into pairs):")
	names := []string{"Alice", "Bob", "Charlie"}
	ages := []int{30, 25, 35}
	pairs := Zip(names, ages)
	for _, p := range pairs {
		fmt.Printf("   %s is %d years old\n", p.First, p.Second)
	}

	fmt.Println("--- Collections & Data Structures Demo End ---")
}
