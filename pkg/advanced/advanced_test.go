package advanced

import (
	"context"
	"math"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Advanced Feature Tests
//
// For a Java developer:
// - Go's testing style is often more "table-driven" (see TestSafeDivide for example).
// - There is no `@Before` or `@After` annotation by default; use `t.Cleanup()`
//   or just defer at the beginning of the test.
// - Test state is managed locally within each test function.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Embedding
// ---------------------------------------------------------------------------

func TestDogEmbedding(t *testing.T) {
	dog := Dog{Animal: Animal{Name: "Rex", Legs: 4}, Breed: "Labrador"}

	// Promoted fields from embedded Animal
	if dog.Name != "Rex" {
		t.Errorf("Expected Name 'Rex', got %s", dog.Name)
	}
	if dog.Legs != 4 {
		t.Errorf("Expected Legs 4, got %d", dog.Legs)
	}

	// Promoted method from embedded Animal
	desc := dog.Describe()
	if !strings.Contains(desc, "Rex") {
		t.Errorf("Describe() should contain 'Rex', got %s", desc)
	}

	// Subclass-like method on Dog
	speak := dog.Speak()
	if !strings.Contains(speak, "Woof") {
		t.Errorf("Speak() should contain 'Woof', got %s", speak)
	}
}

func TestServiceDogOverride(t *testing.T) {
	sd := ServiceDog{
		Dog:    Dog{Animal: Animal{Name: "Buddy", Legs: 4}, Breed: "Golden"},
		CertID: "SD-99",
	}

	// Overridden Describe() method
	desc := sd.Describe()
	if !strings.Contains(desc, "Service Dog") || !strings.Contains(desc, "SD-99") {
		t.Errorf("ServiceDog.Describe() should mention cert, got %s", desc)
	}
}

// ---------------------------------------------------------------------------
// Enum (iota)
// ---------------------------------------------------------------------------

func TestColorEnum(t *testing.T) {
	if Red.String() != "Red" {
		t.Errorf("Expected 'Red', got %s", Red.String())
	}
	if !Red.IsWarm() {
		t.Error("Red should be warm")
	}
	if Blue.IsWarm() {
		t.Error("Blue should not be warm")
	}
	all := AllColors()
	if len(all) != 4 {
		t.Errorf("Expected 4 colors, got %d", len(all))
	}
}

// ---------------------------------------------------------------------------
// Functional patterns
// ---------------------------------------------------------------------------

func TestCompose(t *testing.T) {
	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }
	// Compose(addOne, double) evaluates addOne(double(3)) = 3*2 + 1 = 7
	fn := Compose(addOne, double)
	if fn(3) != 7 {
		t.Errorf("Expected 7, got %d", fn(3))
	}
}

func TestCurry2(t *testing.T) {
	add := func(a, b int) int { return a + b }
	add5 := Curry2(add)(5)
	if add5(3) != 8 {
		t.Errorf("Expected 8, got %d", add5(3))
	}
}

func TestMemoize(t *testing.T) {
	calls := 0
	fn := Memoize(func(n int) int {
		calls++
		return n * n
	})

	fn(3)
	fn(3)
	fn(3)
	if calls != 1 {
		t.Errorf("Expected 1 computation, got %d", calls)
	}
	if fn(4) != 16 {
		t.Errorf("Expected 16, got %d", fn(4))
	}
	if calls != 2 {
		t.Errorf("Expected 2 computations after new key, got %d", calls)
	}
}

func TestMemoizeConcurrentCalls(t *testing.T) {
	var calls atomic.Int64
	memoized := Memoize(func(n int) int {
		calls.Add(1)
		return n * n
	})

	const goroutines = 50
	var wg sync.WaitGroup
	results := make([]int, goroutines)
	wg.Add(goroutines)
	for i := range goroutines {
		go func() {
			defer wg.Done()
			results[i] = memoized(i % 10)
		}()
	}
	wg.Wait()

	for i, result := range results {
		if want := (i % 10) * (i % 10); result != want {
			t.Errorf("memoized(%d) = %d, want %d", i%10, result, want)
		}
	}
	if calls.Load() < 10 {
		t.Errorf("function called %d times, want at least one call per key", calls.Load())
	}
}

func TestPipeline(t *testing.T) {
	result := Pipeline("  hello  ",
		strings.TrimSpace,
		strings.ToUpper,
	)
	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got %s", result)
	}
}

// ---------------------------------------------------------------------------
// Defer / Panic / Recover
// ---------------------------------------------------------------------------

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		want      int
		wantError bool
	}{
		{"normal division", 10, 2, 5, false},
		{"divide by zero", 10, 0, 0, true},
		{"negative result", -10, 2, -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeDivide(tt.a, tt.b)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeDivide() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("SafeDivide() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Reflection
// ---------------------------------------------------------------------------

func TestDescribeType(t *testing.T) {
	type Sample struct {
		X int
		Y string
	}
	desc := DescribeType(Sample{X: 42, Y: "hello"})
	if !strings.Contains(desc, "Sample") {
		t.Errorf("Expected type name 'Sample' in description, got %s", desc)
	}
	if !strings.Contains(desc, "Fields: 2") {
		t.Errorf("Expected 'Fields: 2', got %s", desc)
	}
}

func TestStructToMap(t *testing.T) {
	type Sample struct {
		Name string
		Age  int
	}
	m := StructToMap(Sample{Name: "Alice", Age: 30})
	if m["Name"] != "Alice" {
		t.Errorf("Expected Name='Alice', got %v", m["Name"])
	}
	if m["Age"] != 30 {
		t.Errorf("Expected Age=30, got %v", m["Age"])
	}
}

// ---------------------------------------------------------------------------
// Type constraints (Number)
// ---------------------------------------------------------------------------

func TestSum(t *testing.T) {
	if Sum([]int{1, 2, 3, 4}) != 10 {
		t.Errorf("Expected 10, got %d", Sum([]int{1, 2, 3, 4}))
	}
	if math.Abs(Sum([]float64{1.5, 2.5})-4.0) > 1e-9 {
		t.Errorf("Expected 4.0, got %f", Sum([]float64{1.5, 2.5}))
	}
}

func TestMinMax(t *testing.T) {
	items := []int{3, 1, 4, 1, 5, 9}
	if Min(items) != 1 {
		t.Errorf("Expected Min=1, got %d", Min(items))
	}
	if Max(items) != 9 {
		t.Errorf("Expected Max=9, got %d", Max(items))
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		val, lo, hi int
		want        int
	}{
		{15, 0, 10, 10},
		{-5, 0, 10, 0},
		{5, 0, 10, 5},
	}
	for _, tt := range tests {
		got := Clamp(tt.val, tt.lo, tt.hi)
		if got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d; want %d", tt.val, tt.lo, tt.hi, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// FanOut
// ---------------------------------------------------------------------------

func TestFanOut(t *testing.T) {
	inputs := []int{1, 2, 3, 4, 5}
	results, err := FanOut(inputs, 2, func(n int) int { return n * n })
	if err != nil {
		t.Fatalf("FanOut returned error: %v", err)
	}
	expected := []int{1, 4, 9, 16, 25}
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Expected %v, got %v", expected, results)
	}
}

func TestFanOutRejectsInvalidArguments(t *testing.T) {
	if _, err := FanOut([]int{1}, 0, func(n int) int { return n }); err == nil {
		t.Error("FanOut should reject a non-positive worker count")
	}
	if _, err := FanOut[int, int]([]int{1}, 1, nil); err == nil {
		t.Error("FanOut should reject a nil function")
	}
}

func TestFanOutRespectsWorkerLimit(t *testing.T) {
	var mu sync.Mutex
	active, peak := 0, 0
	_, err := FanOut(make([]int, 30), 3, func(int) int {
		mu.Lock()
		active++
		if active > peak {
			peak = active
		}
		mu.Unlock()

		time.Sleep(time.Millisecond)

		mu.Lock()
		active--
		mu.Unlock()
		return 0
	})
	if err != nil {
		t.Fatalf("FanOut returned error: %v", err)
	}
	if peak > 3 {
		t.Fatalf("peak concurrent calls = %d, want at most 3", peak)
	}
}

// ---------------------------------------------------------------------------
// Java FAQ Tests
// ---------------------------------------------------------------------------

func TestAccount(t *testing.T) {
	acc := Account{Owner: "Alice"}
	acc.Deposit(100)
	if acc.GetBalance() != 100 {
		t.Errorf("Expected balance 100, got %d", acc.GetBalance())
	}
	acc.Deposit(-50)
	if acc.GetBalance() != 100 {
		t.Errorf("Expected balance 100 after negative deposit, got %d", acc.GetBalance())
	}
}

func TestUserRecord(t *testing.T) {
	u := UserRecord{Username: "gopher", Email: "go@golang.org"}
	if u.Username != "gopher" {
		t.Errorf("Expected Username 'gopher', got %s", u.Username)
	}
}

// ---------------------------------------------------------------------------
// Advanced Samples Tests
// ---------------------------------------------------------------------------

func TestDatabaseOptions(t *testing.T) {
	db := NewDatabaseConnector(
		WithHost("example.com"),
		WithPort(9999),
		WithTimeout(5*time.Second),
	)
	if db.Host != "example.com" {
		t.Errorf("Expected Host 'example.com', got %s", db.Host)
	}
	if db.Port != 9999 {
		t.Errorf("Expected Port 9999, got %d", db.Port)
	}
	if db.Timeout != 5*time.Second {
		t.Errorf("Expected Timeout 5s, got %v", db.Timeout)
	}
}

func TestSafeCounter(t *testing.T) {
	counter := &SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	if counter.Value() != 1000 {
		t.Errorf("Expected Value 1000, got %d", counter.Value())
	}
}

func TestConcurrentCache(t *testing.T) {
	cache := NewConcurrentCache[string, string]()
	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Expected value1, got %v, ok: %v", val, ok)
	}
	_, ok = cache.Get("key2")
	if ok {
		t.Error("Expected key2 not to be in cache")
	}
}

func TestAdvancedPipeline(t *testing.T) {
	// Squares: 1, 4, 9, 16. Evens: 4, 16. Sum: 20.
	ctx := context.Background()
	sum, err := AdvancedPipeline(ctx, []int{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sum != 20 {
		t.Errorf("Expected sum 20, got %d", sum)
	}

	// Verify cancellation with a pre-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = AdvancedPipeline(ctx, []int{1, 2, 3})
	if err == nil {
		t.Error("Expected error from cancelled context, got nil")
	}
}
