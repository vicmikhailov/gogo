package collections

import (
	"reflect"
	"strings"
	"testing"
)

/**
 * ===========================================================================
 * Collections Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - `t.Run` allows for subtests, similar to JUnit 5 `@Nested` tests.
 * - This keeps tests organized when testing different aspects of a structure.
 */

// ---------------------------------------------------------------------------
// Set tests
// ---------------------------------------------------------------------------

func TestSetBasicOps(t *testing.T) {
	s := NewSet(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Expected Len 3, got %d", s.Len())
	}
	if !s.Contains(2) {
		t.Errorf("Expected Contains(2) = true")
	}
	s.Remove(2)
	if s.Contains(2) {
		t.Errorf("Expected Contains(2) = false after Remove")
	}
	if s.Len() != 2 {
		t.Errorf("Expected Len 2 after Remove, got %d", s.Len())
	}
}

func TestSetUnion(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(3, 4, 5)
	u := a.Union(b)
	if u.Len() != 5 {
		t.Errorf("Expected union Len 5, got %d", u.Len())
	}
	for _, v := range []int{1, 2, 3, 4, 5} {
		if !u.Contains(v) {
			t.Errorf("Expected union to contain %d", v)
		}
	}
}

func TestSetIntersection(t *testing.T) {
	a := NewSet(1, 2, 3, 4)
	b := NewSet(3, 4, 5, 6)
	inter := a.Intersection(b)
	got := Sorted(inter.Values(), func(x, y int) bool { return x < y })
	expected := []int{3, 4}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected intersection %v, got %v", expected, got)
	}
}

func TestSetDifference(t *testing.T) {
	a := NewSet(1, 2, 3, 4)
	b := NewSet(3, 4, 5, 6)
	diff := a.Difference(b)
	got := Sorted(diff.Values(), func(x, y int) bool { return x < y })
	expected := []int{1, 2}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected difference %v, got %v", expected, got)
	}
}

// ---------------------------------------------------------------------------
// Stack tests
// ---------------------------------------------------------------------------

func TestStack(t *testing.T) {
	s := Stack[int]{}
	if !s.IsEmpty() {
		t.Error("New stack should be empty")
	}
	s.Push(10)
	s.Push(20)
	s.Push(30)

	if s.Len() != 3 {
		t.Errorf("Expected Len 3, got %d", s.Len())
	}
	val, ok := s.Peek()
	if !ok || val != 30 {
		t.Errorf("Expected Peek 30, got %d", val)
	}

	val, ok = s.Pop()
	if !ok || val != 30 {
		t.Errorf("Expected Pop 30, got %d", val)
	}
	val, ok = s.Pop()
	if !ok || val != 20 {
		t.Errorf("Expected Pop 20, got %d", val)
	}
	val, ok = s.Pop()
	if !ok || val != 10 {
		t.Errorf("Expected Pop 10, got %d", val)
	}

	_, ok = s.Pop()
	if ok {
		t.Error("Pop on empty stack should return ok=false")
	}
}

// ---------------------------------------------------------------------------
// Queue tests
// ---------------------------------------------------------------------------

func TestQueue(t *testing.T) {
	q := Queue[string]{}
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	q.Enqueue("a")
	q.Enqueue("b")
	q.Enqueue("c")

	if q.Len() != 3 {
		t.Errorf("Expected Len 3, got %d", q.Len())
	}

	val, ok := q.Dequeue()
	if !ok || val != "a" {
		t.Errorf("Expected Dequeue 'a', got %s", val)
	}
	val, ok = q.Dequeue()
	if !ok || val != "b" {
		t.Errorf("Expected Dequeue 'b', got %s", val)
	}
	val, ok = q.Dequeue()
	if !ok || val != "c" {
		t.Errorf("Expected Dequeue 'c', got %s", val)
	}
	_, ok = q.Dequeue()
	if ok {
		t.Error("Dequeue on empty queue should return ok=false")
	}
}

// ---------------------------------------------------------------------------
// OrderedMap tests
// ---------------------------------------------------------------------------

func TestOrderedMapPreservesInsertionOrder(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Put("charlie", 3)
	om.Put("alpha", 1)
	om.Put("bravo", 2)

	expectedKeys := []string{"charlie", "alpha", "bravo"}
	if !reflect.DeepEqual(om.Keys(), expectedKeys) {
		t.Errorf("Expected keys %v, got %v", expectedKeys, om.Keys())
	}

	v, ok := om.Get("alpha")
	if !ok || v != 1 {
		t.Errorf("Expected Get('alpha') = 1, got %d", v)
	}

	om.Delete("alpha")
	expectedKeys = []string{"charlie", "bravo"}
	if !reflect.DeepEqual(om.Keys(), expectedKeys) {
		t.Errorf("After delete expected keys %v, got %v", expectedKeys, om.Keys())
	}
	if om.Len() != 2 {
		t.Errorf("Expected Len 2, got %d", om.Len())
	}
}

// ---------------------------------------------------------------------------
// Functional operations tests
// ---------------------------------------------------------------------------

func TestFilter(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(items, func(n int) bool { return n%2 == 0 })
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(evens, expected) {
		t.Errorf("Expected %v, got %v", expected, evens)
	}
}

func TestReduce(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	sum := Reduce(items, 0, func(acc, n int) int { return acc + n })
	if sum != 15 {
		t.Errorf("Expected sum 15, got %d", sum)
	}
}

func TestFlatMap(t *testing.T) {
	// String splitting
	words := FlatMap([]string{"hello world", "go lang"}, func(s string) []string {
		return strings.Split(s, " ")
	})
	expectedWords := []string{"hello", "world", "go", "lang"}
	if !reflect.DeepEqual(words, expectedWords) {
		t.Errorf("Expected %v, got %v", expectedWords, words)
	}

	// Flattening nested slices
	flat := FlatMap([][]int{{1, 2}, {3, 4}, {5}}, func(s []int) []int { return s })
	expectedFlat := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(flat, expectedFlat) {
		t.Errorf("Expected %v, got %v", expectedFlat, flat)
	}
}

func TestGroupBy(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	grouped := GroupBy(items, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if !reflect.DeepEqual(grouped["even"], []int{2, 4, 6}) {
		t.Errorf("Expected even group [2,4,6], got %v", grouped["even"])
	}
	if !reflect.DeepEqual(grouped["odd"], []int{1, 3, 5}) {
		t.Errorf("Expected odd group [1,3,5], got %v", grouped["odd"])
	}
}

func TestPartition(t *testing.T) {
	items := []int{-2, -1, 0, 1, 2}
	pos, neg := Partition(items, func(n int) bool { return n >= 0 })
	expectedPos := []int{0, 1, 2}
	expectedNeg := []int{-2, -1}
	if !reflect.DeepEqual(pos, expectedPos) {
		t.Errorf("Expected pos %v, got %v", expectedPos, pos)
	}
	if !reflect.DeepEqual(neg, expectedNeg) {
		t.Errorf("Expected neg %v, got %v", expectedNeg, neg)
	}
}

func TestDistinct(t *testing.T) {
	items := []int{1, 2, 2, 3, 3, 3, 4}
	got := Distinct(items)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestAnyAll(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	if Any(items, func(n int) bool { return n < 0 }) {
		t.Error("Expected Any(negative) = false")
	}
	if !Any(items, func(n int) bool { return n == 3 }) {
		t.Error("Expected Any(==3) = true")
	}
	if !All(items, func(n int) bool { return n > 0 }) {
		t.Error("Expected All(>0) = true")
	}
	if All(items, func(n int) bool { return n > 3 }) {
		t.Error("Expected All(>3) = false")
	}
}

func TestZip(t *testing.T) {
	names := []string{"Alice", "Bob"}
	ages := []int{30, 25, 99} // extra element ignored
	pairs := Zip(names, ages)
	if len(pairs) != 2 {
		t.Errorf("Expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].First != "Alice" || pairs[0].Second != 30 {
		t.Errorf("Expected (Alice,30), got (%s,%d)", pairs[0].First, pairs[0].Second)
	}
	if pairs[1].First != "Bob" || pairs[1].Second != 25 {
		t.Errorf("Expected (Bob,25), got (%s,%d)", pairs[1].First, pairs[1].Second)
	}
}

func TestSorted(t *testing.T) {
	items := []int{5, 3, 1, 4, 2}
	sorted := Sorted(items, func(a, b int) bool { return a < b })
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected %v, got %v", expected, sorted)
	}
	// Original should be unchanged
	if items[0] != 5 {
		t.Error("Sorted should not mutate original slice")
	}
}
