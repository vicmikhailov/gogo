// Package advanced provides a deep dive into Go's advanced features,
// contrasting them with equivalent concepts in Java.
//
// Feature          | Go approach                          | Java approach
// -----------------|--------------------------------------|--------------------------------------
// Objects          | Structs + Methods (Composition)      | Classes + Inheritance
// Interfaces       | Implicit ("Duck Typing")             | Explicit (`implements`)
// Inheritance      | Struct Embedding (Composition)       | Class Inheritance (`extends`)
// Encapsulation    | Exported vs Unexported (Casing)      | `public`, `private`, `protected`
// Error Handling   | Explicit return values (`error`)     | Exceptions (`try-catch-finally`)
// Concurrency      | Goroutines + Channels (CSP)          | Threads, OS processes, Virtual Threads
// Generics         | Compile-time monomorphization        | Type Erasure
// Meta-programming | Reflection, Struct Tags, Code Gen    | Reflection, Annotations, Proxying
// Resource Mgmt    | `defer`                              | `try-with-resources` or `finally`
// Pointers         | Optional (explicit control)          | Implicit (all objects are references)
// Collections      | Slices, Maps (Built-in)              | List, Map (Library-based)
// Null Safety      | Nil (specific types)                 | Null (any object)
//
// Common Java-isms to avoid in Go:
//   - "New" everything: Don't feel forced to create `NewXYZ` if the zero-value is useful.
//   - Getters/Setters: Direct field access is preferred unless logic is needed.
//   - Deep Nesting: Go prefers a flatter package structure.
//   - Over-Interfaces: Don't define an interface before you have at least two implementations.
//   - Panic-as-Exception: Never use `panic` for normal error flow; use `error` returns.
//   - Pointer-to-Everything: Use values by default; only use pointers when mutation or size matters.
//
// Special Agreements & Tooling (Go's unique mechanisms):
//  1. `//go:generate`: Tool-driven code generation (like Annotation Processors).
//  2. Build Tags (`//go:build`): Conditional compilation at file level (no C++ style #define).
//  3. `internal` packages: Strong encapsulation (only siblings/parent can import).
//  4. `init()` functions: Automatic per-package setup (like static initializer blocks).
//  5. `iota`: Auto-incrementing counter for constants and bitmask flags.
//  6. Struct Tags: Compile-time metadata for reflection (like Jackson annotations).
//  7. Records: Go `struct` is the closest equivalent to a Java `record` (data-focused).
//  8. Bean Agreement: Direct field access preferred over Getters/Setters unless validation is required.
package advanced

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//go:embed static/hello.txt
var embeddedFile string

// ---------------------------------------------------------------------------
// 8. Struct Tags and JSON (essential for web development)
// ---------------------------------------------------------------------------

// PersonWithTags demonstrates struct metadata tags for JSON serialization.
//
// For a Java developer:
//   - Struct tags are like Annotations in Java (e.g. Jackson `@JsonProperty("first_name")`).
//   - They provide metadata readable at runtime via reflection.
//   - `omitempty` skips fields with zero values during marshaling.
//   - `json:"-"` ignores the field entirely (like `@JsonIgnore`).
type PersonWithTags struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age,omitempty"`
	Secret    string `json:"-"`
}

// ---------------------------------------------------------------------------
// 1. Embedding (Go's composition – comparable to Java inheritance)
// ---------------------------------------------------------------------------

// Animal is a base data struct demonstrating composition.
//
// For a Java developer:
//   - Go does not have `public`, `private`, or `protected`. Visibility is based on casing:
//   - Uppercase (e.g., `Describe`): Exported / Public (accessible outside package).
//   - Lowercase (e.g., `secret`): Unexported / Private (accessible only within package).
//   - Go does not have classes or class inheritance (`extends`). Instead, it uses composition.
type Animal struct {
	Name string
	Legs int
}

// Describe returns a string description of the animal.
// Equivalent to a method on a Java base class.
func (a Animal) Describe() string {
	return fmt.Sprintf("%s has %d legs", a.Name, a.Legs)
}

// Dog embeds Animal.
//
// For a Java developer:
//   - Anonymous field `Animal` is called "Struct Embedding".
//   - Fields and methods of `Animal` are "promoted" to `Dog`.
//   - Java equivalent: `class Dog extends Animal`.
type Dog struct {
	Animal
	Breed string
}

// Speak returns the sound the dog makes (subclass-like method).
func (d Dog) Speak() string {
	return fmt.Sprintf("%s says: Woof!", d.Name)
}

// ServiceDog embeds Dog and overrides Describe.
type ServiceDog struct {
	Dog
	CertID string
}

// Describe overrides the embedded Animal.Describe.
//
// For a Java developer:
//   - Similar to `@Override` in Java.
//   - Calling `sd.Describe()` invokes this method.
//   - You can still access the embedded "parent" method via `sd.Dog.Describe()`.
func (sd ServiceDog) Describe() string {
	return fmt.Sprintf("%s [Service Dog, Cert: %s]", sd.Dog.Describe(), sd.CertID)
}

// ---------------------------------------------------------------------------
// 2. Enum pattern with iota (comparable to Java enum)
// ---------------------------------------------------------------------------

// Color acts as a type-safe enum using iota.
//
// For a Java developer:
//   - Go doesn't have an `enum` keyword.
//   - We define `type Color int` and a `const` block with `iota`.
//   - `iota` auto-increments by 1 on each line within the const block (starting from 0).
//   - Similar to `public static final int` but with dedicated type safety.
type Color int

const (
	Red Color = iota
	Green
	Blue
	Yellow
)

// String implements the `fmt.Stringer` interface (like Java's `toString()` on an Enum).
func (c Color) String() string {
	return [...]string{"Red", "Green", "Blue", "Yellow"}[c]
}

// IsWarm reports whether the color is considered warm.
func (c Color) IsWarm() bool {
	return c == Red || c == Yellow
}

// AllColors returns all valid Color values (like Java's `Color.values()`).
func AllColors() []Color {
	return []Color{Red, Green, Blue, Yellow}
}

// ---------------------------------------------------------------------------
// 3. Functional programming patterns (closures, currying, memoization)
// ---------------------------------------------------------------------------

// Compose returns the composition of two functions: f(g(x)).
// Java equivalent: `Function.compose()`
func Compose[T any](f, g func(T) T) func(T) T {
	return func(x T) T { return f(g(x)) }
}

// Curry2 converts a two-argument function into a chain of two unary functions.
func Curry2[A, B, C any](fn func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C {
			return fn(a, b)
		}
	}
}

// Memoize wraps a function with thread-safe caching of previous results.
// Java equivalent: Caching results in a `ConcurrentHashMap` or using Guava's Cache.
// Concurrent misses for the same key may compute fn more than once; fn runs outside the lock.
func Memoize[K comparable, V any](fn func(K) V) func(K) V {
	cache := make(map[K]V)
	var mu sync.Mutex
	return func(key K) V {
		mu.Lock()
		if val, ok := cache[key]; ok {
			mu.Unlock()
			return val
		}
		mu.Unlock()

		val := fn(key)

		mu.Lock()
		defer mu.Unlock()
		if cached, ok := cache[key]; ok {
			return cached
		}
		cache[key] = val
		return val
	}
}

// Pipeline applies a sequence of functions to an initial value from left to right.
// Java equivalent: Chaining multiple `Function.andThen()` calls.
func Pipeline[T any](value T, fns ...func(T) T) T {
	for _, fn := range fns {
		value = fn(value)
	}
	return value
}

// ---------------------------------------------------------------------------
// 4. Defer / Panic / Recover (comparable to Java try-catch-finally)
// ---------------------------------------------------------------------------

// SafeDivide demonstrates catching a panic and returning an error.
//
// For a Java developer:
//   - Go idiom: `panic` is for truly unrecoverable errors (like `OutOfMemoryError` or programmer bugs).
//   - `recover()` can only be called inside a `defer` block to catch an active panic.
//   - Java equivalent: `try-catch` on `ArithmeticException`.
func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}()
	return a / b, nil
}

// WithResource demonstrates deterministic cleanup via `defer`.
//
// For a Java developer:
//   - Go idiom: `defer` is the standard way to guarantee cleanup (closing files, releasing locks).
//   - Multiple `defer` calls execute in LIFO (Last-In-First-Out) order.
//   - Java equivalent: `try-with-resources` or `finally` block.
func WithResource(name string) []string {
	var log []string
	log = append(log, fmt.Sprintf("Opening resource: %s", name))
	defer func() {
		log = append(log, fmt.Sprintf("Closing resource: %s (via defer)", name))
	}()
	log = append(log, fmt.Sprintf("Using resource: %s", name))
	return log
}

// ---------------------------------------------------------------------------
// 5. Reflection (comparable to Java reflection API)
// ---------------------------------------------------------------------------

// DescribeType inspects a value's type, kind, and fields at runtime.
//
// For a Java developer:
//   - Java equivalent: `object.getClass()`, `class.getDeclaredFields()`.
//   - Go's `reflect` package provides `reflect.TypeOf` and `reflect.ValueOf`.
//   - Should be used sparingly in Go due to runtime overhead and loss of type safety.
func DescribeType(v interface{}) string {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Type: %s, Kind: %s", t.Name(), t.Kind()))

	if t.Kind() == reflect.Struct {
		sb.WriteString(fmt.Sprintf(", Fields: %d", t.NumField()))
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldVal := val.Field(i)
			sb.WriteString(fmt.Sprintf("\n  - %s (%s) = %v", field.Name, field.Type, fieldVal))
		}
	}
	return sb.String()
}

// StructToMap converts any struct's exported fields into a map using reflection.
// Java equivalent: Jackson's `objectMapper.convertValue(pojo, Map.class)`.
func StructToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return result
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			result[field.Name] = val.Field(i).Interface()
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// 6. Type constraints and type sets (advanced generics)
// ---------------------------------------------------------------------------

// Number constrains generic types to all built-in numeric types.
//
// For a Java developer:
//   - Similar to `<T extends Number>` in Java generics.
//   - The `~` symbol allows custom types whose underlying type is one of these primitives.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum calculates the sum of all elements in a numeric slice.
// Java equivalent: `public <T extends Number> T sum(List<T> items)`
func Sum[T Number](items []T) T {
	var total T
	for _, v := range items {
		total += v
	}
	return total
}

// Min returns the minimum element in a numeric slice.
// Java equivalent: `Collections.min()`
func Min[T Number](items []T) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	m := items[0]
	for _, v := range items[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// Max returns the maximum element in a numeric slice.
// Java equivalent: `Collections.max()`
func Max[T Number](items []T) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	m := items[0]
	for _, v := range items[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// Clamp restricts value to the range [lo, hi].
// Java equivalent: `Math.clamp()` (introduced in Java 21).
func Clamp[T Number](value, lo, hi T) T {
	if value < lo {
		return lo
	}
	if value > hi {
		return hi
	}
	return value
}

// ---------------------------------------------------------------------------
// 7. Concurrent patterns (advanced – comparable to java.util.concurrent)
// ---------------------------------------------------------------------------

// FanOut distributes work across a fixed pool of concurrent workers and preserves result ordering.
// It returns an error if workers is not positive or fn is nil.
//
// For a Java developer:
//   - Java equivalent: `ExecutorService.invokeAll()`, parallel streams, or `CompletableFuture.allOf()`.
//   - Go idiom: A jobs channel feeds a bounded number of worker goroutines, and
//     `sync.WaitGroup` waits for every worker to complete.
func FanOut[T any, R any](items []T, workers int, fn func(T) R) ([]R, error) {
	if workers <= 0 {
		return nil, fmt.Errorf("worker count must be positive")
	}
	if fn == nil {
		return nil, fmt.Errorf("fan-out function must not be nil")
	}
	if len(items) == 0 {
		return []R{}, nil
	}
	if workers > len(items) {
		workers = len(items)
	}

	type indexed struct {
		idx   int
		value T
	}
	type indexedResult struct {
		idx    int
		result R
	}
	jobs := make(chan indexed, len(items))
	results := make(chan indexedResult, len(items))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- indexedResult{idx: job.idx, result: fn(job.value)}
			}
		}()
	}

	for i, item := range items {
		jobs <- indexed{idx: i, value: item}
	}
	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()

	ordered := make([]R, len(items))
	for result := range results {
		ordered[result.idx] = result.result
	}
	return ordered, nil
}

// ---------------------------------------------------------------------------
// 11. Functional Options Pattern (Common Go pattern for configuration)
// ---------------------------------------------------------------------------

// DatabaseConfig holds configuration for a database client.
//
// For a Java developer:
//   - This is an idiomatic Go alternative to the Builder pattern.
//   - It avoids telescoping constructors and handles defaults cleanly.
type DatabaseConfig struct {
	Host     string
	Port     int
	Timeout  time.Duration
	MaxConns int
}

// DBOption modifies DatabaseConfig.
type DBOption func(*DatabaseConfig)

// WithHost sets the host option (Java: builder.setHost(host)).
func WithHost(host string) DBOption {
	return func(c *DatabaseConfig) {
		c.Host = host
	}
}

// WithPort sets the port option (Java: builder.setPort(port)).
func WithPort(port int) DBOption {
	return func(c *DatabaseConfig) {
		c.Port = port
	}
}

// WithTimeout sets the timeout option (Java: builder.setTimeout(t)).
func WithTimeout(t time.Duration) DBOption {
	return func(c *DatabaseConfig) {
		c.Timeout = t
	}
}

// NewDatabaseConnector constructs a DatabaseConfig applying functional options over defaults.
func NewDatabaseConnector(opts ...DBOption) *DatabaseConfig {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Timeout:  30 * time.Second,
		MaxConns: 10,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// ---------------------------------------------------------------------------
// 12. Atomic Operations (Low-level synchronization)
// ---------------------------------------------------------------------------

// SafeCounter demonstrates thread-safe lock-free counting via `sync/atomic`.
//
// For a Java developer:
//   - Equivalent to `java.util.concurrent.atomic.AtomicLong`.
//   - More performant than `sync.Mutex` for simple counter increments.
type SafeCounter struct {
	value atomic.Int64
}

// Increment safely increments the counter (Java: `atomicLong.incrementAndGet()`).
func (c *SafeCounter) Increment() {
	c.value.Add(1)
}

// Value returns the current value (Java: `atomicLong.get()`).
func (c *SafeCounter) Value() int64 {
	return c.value.Load()
}

// ---------------------------------------------------------------------------
// 13. Compile-time Interface Satisfaction Check
// ---------------------------------------------------------------------------

// Logger defines a simple logging contract.
type Logger interface {
	Log(message string)
}

// ConsoleLogger implements the Logger interface.
type ConsoleLogger struct{}

func (c ConsoleLogger) Log(message string) {
	fmt.Println("   [Logger]: " + message)
}

// Compile-time check ensuring ConsoleLogger implements Logger.
//
// For a Java developer:
//   - Java validates interface satisfaction at compile-time via `class Foo implements Bar`.
//   - In Go (where interfaces are implicit), this idiom forces a compile error if ConsoleLogger
//     fails to satisfy Logger: `var _ Interface = (*Implementation)(nil)`.
var _ Logger = (*ConsoleLogger)(nil)

// ---------------------------------------------------------------------------
// 14. Generic Concurrent Cache (RWMutex + Generics)
// ---------------------------------------------------------------------------

// ConcurrentCache provides thread-safe key-value caching using `sync.RWMutex`.
//
// For a Java developer:
//   - Similar to `ConcurrentHashMap<K, V>`.
//   - `sync.RWMutex` allows multiple concurrent readers but only one writer (like `ReentrantReadWriteLock`).
type ConcurrentCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// NewConcurrentCache initializes a generic thread-safe cache.
func NewConcurrentCache[K comparable, V any]() *ConcurrentCache[K, V] {
	return &ConcurrentCache[K, V]{
		items: make(map[K]V),
	}
}

// Set adds or updates a value in the cache under write lock.
func (c *ConcurrentCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// Get retrieves a value from the cache under read lock.
func (c *ConcurrentCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

// ---------------------------------------------------------------------------
// 15. Context-Aware Multi-stage Pipeline (Advanced Concurrency)
// ---------------------------------------------------------------------------

type pipelineResult struct {
	val int
	err error
}

// AdvancedPipeline executes a multi-stage concurrency pipeline with cancellation support.
//
// For a Java developer:
//   - Java equivalent: CompletableFuture chain or Reactive Streams with cancellation.
//   - Demonstrates multi-stage processing using channels and context for graceful cancellation.
func AdvancedPipeline(ctx context.Context, nums []int) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	// Stage 1: Square numbers and stream through channel
	stage1 := make(chan pipelineResult)
	go func() {
		defer close(stage1)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case stage1 <- pipelineResult{val: n * n}:
			}
		}
	}()

	// Stage 2: Filter even squares and sum them
	sum := 0
	for res := range stage1 {
		if res.err != nil {
			return 0, res.err
		}
		if res.val%2 == 0 {
			sum += res.val
		}

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
	}

	return sum, nil
}

// ---------------------------------------------------------------------------
// 16. Special Agreements (Tooling)
// ---------------------------------------------------------------------------

//go:generate echo "Running custom code generation tool..."

// Account demonstrates field visibility and idiomatic getter/setter usage in Go.
//
// For a Java developer:
//   - Go does not use the JavaBeans pattern (get/set for every field).
//   - Exported fields (capitalized) are preferred for direct access.
//   - Getters/setters are reserved for logic (e.g. validation, mutation of unexported state).
type Account struct {
	id      int
	balance int
	Owner   string
}

// GetBalance is a getter returning balance using a value receiver.
func (a Account) GetBalance() int {
	return a.balance
}

// Deposit is a setter mutating balance using a pointer receiver.
func (a *Account) Deposit(amount int) {
	if amount > 0 {
		a.balance += amount
	}
}

// UserRecord represents a lightweight data-holder struct (comparable to Java 14+ record).
type UserRecord struct {
	Username string
	Email    string
}

// RunAdvancedDemo orchestrates all advanced language feature demonstrations.
func RunAdvancedDemo() {
	fmt.Println("--- Advanced Language Features Demo ---")

	// 1. Embedding (composition ≈ Java inheritance)
	fmt.Println("1. Embedding (composition ≈ Java inheritance):")
	dog := Dog{Animal: Animal{Name: "Rex", Legs: 4}, Breed: "German Shepherd"}
	fmt.Printf("   %s (Breed: %s)\n", dog.Describe(), dog.Breed)
	fmt.Printf("   %s\n", dog.Speak())
	sd := ServiceDog{Dog: dog, CertID: "SD-42"}
	fmt.Printf("   %s\n", sd.Describe())

	// 2. Enums with iota
	fmt.Println("2. Enum pattern (iota ≈ Java enum):")
	for _, c := range AllColors() {
		fmt.Printf("   %s (ordinal=%d, warm=%v)\n", c, c, c.IsWarm())
	}

	// 3. Functional patterns (Compose, Curry, Memoize, Pipeline)
	fmt.Println("3. Functional programming:")
	double := func(x int) int { return x * 2 }
	addThree := func(x int) int { return x + 3 }
	doubleThenAdd := Compose(addThree, double)
	fmt.Printf("   Compose(+3, *2)(5) = %d\n", doubleThenAdd(5))

	curriedAdd := Curry2(func(a, b int) int { return a + b })
	add10 := curriedAdd(10)
	fmt.Printf("   Curry add(10)(5) = %d\n", add10(5))

	calls := 0
	expensive := Memoize(func(n int) int {
		calls++
		return n * n
	})
	expensive(4)
	expensive(4)
	expensive(4)
	fmt.Printf("   Memoize: square(4) called 3 times, computed %d time(s)\n", calls)

	result := Pipeline("  Hello, World!  ",
		strings.TrimSpace,
		strings.ToUpper,
		func(s string) string { return ">>>" + s + "<<<" },
	)
	fmt.Printf("   Pipeline: %s\n", result)

	// 4. Defer / Panic / Recover
	fmt.Println("4. Defer/Panic/Recover (≈ try-catch-finally):")
	val, err := SafeDivide(10, 3)
	fmt.Printf("   SafeDivide(10, 3) = %d, err = %v\n", val, err)
	val, err = SafeDivide(10, 0)
	fmt.Printf("   SafeDivide(10, 0) = %d, err = %v\n", val, err)

	fmt.Printf("   Defer ordering: ")
	func() {
		defer fmt.Print("third ")
		defer fmt.Print("second ")
		defer fmt.Print("first ")
	}()
	fmt.Println("(defers execute LIFO)")

	// 5. Reflection
	fmt.Println("5. Reflection (≈ Java reflection API):")
	type User struct {
		Name  string
		Age   int
		Email string
	}
	u := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
	fmt.Printf("   %s\n", DescribeType(u))
	m := StructToMap(u)
	fmt.Printf("   StructToMap: %v\n", m)

	// 6. Type constraints (advanced generics)
	fmt.Println("6. Type constraints (Number interface):")
	ints := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Printf("   Sum(%v) = %d\n", ints, Sum(ints))
	fmt.Printf("   Min = %d, Max = %d\n", Min(ints), Max(ints))
	fmt.Printf("   Clamp(15, 0, 10) = %d\n", Clamp(15, 0, 10))
	floats := []float64{3.14, 2.71, 1.41}
	fmt.Printf("   Sum(%v) = %.2f\n", floats, Sum(floats))

	// 7. Concurrent fan-out
	fmt.Println("7. Concurrent fan-out (≈ ExecutorService):")
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	squares, err := FanOut(inputs, 3, func(n int) int { return n * n })
	if err != nil {
		fmt.Printf("   FanOut failed: %v\n", err)
	} else {
		fmt.Printf("   FanOut squares: %v\n", squares)
	}

	// 8. Struct Tags and JSON
	fmt.Println("8. Struct Tags and JSON (common in REST APIs):")
	p := PersonWithTags{FirstName: "John", LastName: "Doe", Age: 30, Secret: "password123"}
	data, err := json.MarshalIndent(p, "   ", "  ")
	if err != nil {
		fmt.Printf("   Could not marshal person: %v\n", err)
	} else {
		fmt.Printf("   Marshaled: %s\n", string(data))
	}
	fmt.Println("   (Note: 'Secret' is ignored and 'Age' would be omitted if it were zero)")

	// 9. Embed (Go 1.16+)
	fmt.Println("9. Embed package (embedding static assets):")
	fmt.Printf("   Embedded file content:\n   %s", embeddedFile)

	// 10. Advanced Iota (Bitmask flags)
	fmt.Println("10. Advanced Iota (Bitmask flags):")
	type Permission int
	const (
		Read Permission = 1 << iota
		Write
		Execute
	)
	pReadWrite := Read | Write
	fmt.Printf("   Permissions: %d (Read:%v, Write:%v, Execute:%v)\n",
		pReadWrite, pReadWrite&Read != 0, pReadWrite&Write != 0, pReadWrite&Execute != 0)

	// 11. Functional Options Pattern
	fmt.Println("11. Functional Options Pattern (Go configuration):")
	dbConfig := NewDatabaseConnector(
		WithHost("db.example.com"),
		WithPort(5433),
		WithTimeout(10*time.Second),
	)
	fmt.Printf("   DBConfig: %+v\n", dbConfig)

	// 12. Atomic operations (lock-free synchronization)
	fmt.Println("12. Atomic operations (lock-free synchronization):")
	counter := &SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("   SafeCounter value (100 parallel increments): %d\n", counter.Value())

	// 13. Compile-time Interface Check
	fmt.Println("13. Compile-time Interface Check:")
	var _ Logger = (*ConsoleLogger)(nil)
	logger := &ConsoleLogger{}
	logger.Log("Hello from the interface-validated logger!")

	// 14. Generic Concurrent Cache (RWMutex + Generics)
	fmt.Println("14. Generic Concurrent Cache (RWMutex + Generics):")
	cache := NewConcurrentCache[string, int]()
	cache.Set("Go", 2009)
	cache.Set("Java", 1995)
	if val, ok := cache.Get("Go"); ok {
		fmt.Printf("   Cache: Go was released in %d\n", val)
	}

	// 15. Advanced Context-Aware Pipeline
	fmt.Println("15. Context-Aware Multi-stage Pipeline:")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	pipelineSum, err := AdvancedPipeline(ctx, []int{1, 2, 3, 4, 5, 6})
	if err != nil {
		fmt.Printf("   Pipeline error: %v\n", err)
	} else {
		fmt.Printf("   Pipeline Sum of even squares: %d\n", pipelineSum)
	}

	// 16. Java Developer FAQ
	fmt.Println("16. Java Developer FAQ (Getters/Setters, Records, Annotations):")
	RunJavaDeveloperFAQ()

	// 17. Build System Comparison
	fmt.Println("17. Build System Comparison (Go vs Maven):")
	RunBuildSystemDemo()

	fmt.Println("--- Advanced Language Features Demo End ---")
}

// RunJavaDeveloperFAQ provides explicit examples for common Java vs Go questions.
func RunJavaDeveloperFAQ() {
	// a. Getters and Setters (Bean Agreement)
	// Java-ism: Writing GetX() and SetX() for every field.
	// Go way: Use direct field access (Uppercase = Public).
	// Only use methods if you need to encapsulate internal logic or provide a read-only view.
	acc := Account{id: 1, balance: 100, Owner: "Alice"}
	acc.Deposit(50)
	fmt.Printf("   - Bean Agreement: Owner is %s (Direct), Balance is %d (Getter)\n", acc.Owner, acc.GetBalance())

	// b. Interface Definition (Producer vs Consumer)
	// Java-ism: Defining Shape interface in the same package as Rectangle/Circle.
	// Go way: "Accept interfaces, return structs".
	// Interfaces should be defined where they are USED (the consumer), not where they are implemented.
	// This allows the consumer to define only the contract they need.
	fmt.Println("   - Interface Tip: Define interfaces where they are USED, not where the struct is defined.")

	// c. Panic vs Error
	// Java-ism: Using panic() like Throw new RuntimeException().
	// Go way: Always return error. Use panic only for unrecoverable programmer errors (like index out of bounds).
	fmt.Println("   - Errors: Return error, don't throw. Go doesn't have try-catch.")

	// d. Constructors and Zero-Values
	// Java-ism: Thinking everything needs a New...() constructor.
	// Go way: Design structs so their "zero-value" (all fields 0/nil) is actually useful.
	// For example, a sync.Mutex or a bytes.Buffer don't need a New...() call.
	var buf strings.Builder
	buf.WriteString("   - Zero-value: Ready to use without a constructor!")
	fmt.Println(buf.String())

	// e. Pointers vs Values
	// Java-ism: Passing everything as a pointer (*T) because Java objects are references.
	// Go way: Pass by value (T) by default. Use pointers ONLY if you need to:
	// 1. Mutate the original object.
	// 2. Avoid copying a large struct when measurement or API semantics justify it.
	// 3. Represent 'nil' for an optional value.
	fmt.Println("   - Pointers: Use values by default; use pointers only when mutation or size matters.")

	// f. Records vs Structs
	// Records are shallowly immutable; Go structs can expose mutable fields. A value
	// receiver copies the receiver but does not make referenced fields immutable.
	u := UserRecord{"gopher", "go@golang.org"}
	fmt.Printf("   - Record-like Struct: %+v\n", u)

	// g. Annotations vs Struct Tags
	// Java Annotations can be processed at compile-time or runtime.
	// Struct tags are only available via reflection at runtime.
	fmt.Println("   - Annotations: Use Struct Tags for metadata (Runtime) and go:generate (Compile-time).")

	// h. Threads vs Goroutines
	// The runtime schedules goroutines onto OS threads; cost varies with workload and runtime version.
	fmt.Println("   - Threads: Goroutines are lightweight tasks scheduled by the Go runtime onto OS threads.")

	// i. No Magic (AOP, Annotations, Dependency Injection)
	// Java-ism: Relying on reflection/proxies for business logic (e.g. @Transactional).
	// Go way: Be explicit. Pass dependencies as arguments; use closures or decorators for middleware.
	// Go avoids "magic" that happens behind the scenes.
	fmt.Println("   - No Magic: Go prefers explicit code over heavy AOP or reflection-based frameworks.")
}

// RunBuildSystemDemo contrasts the Go toolchain with Maven/Gradle.
func RunBuildSystemDemo() {
	// 1. Build Tool: Go uses the 'go' tool (CLI-first)
	fmt.Println("   1. Build Tool: Go uses the 'go' tool (CLI-first).")
	fmt.Println("      - 'go build': Compiles the project (≈ 'mvn package'). Produces a self-contained binary.")
	fmt.Println("      - 'go test': Runs all tests (≈ 'mvn test'). Native support for unit tests, benchmarks, and fuzzing.")
	fmt.Println("      - 'go run': Compiles and runs (≈ 'mvn exec:java'). No need for manual compilation during dev.")

	// 2. Dependency Management: Go Modules ('go.mod')
	fmt.Println("   2. Dependency Management: Go uses Go Modules ('go.mod').")
	fmt.Println("      - No XML! It's a simple text file (≈ 'pom.xml').")
	fmt.Println("      - 'go mod tidy': Syncs dependencies (≈ Maven Import/Refresh). It also cleans up unused ones.")
	fmt.Println("      - Dependencies are downloaded to $GOPATH/pkg/mod, shared across projects (≈ ~/.m2/repository).")

	// 3. Project Structure and Plugins
	fmt.Println("   3. Plugins & Tooling:")
	fmt.Println("      - In Maven, you use plugins (maven-compiler-plugin, etc.).")
	fmt.Println("      - In Go, the 'go' tool is batteries-included (fmt, vet, test, doc, build).")
	fmt.Println("      - For custom tasks, Go devs use a 'Makefile' (common in C/C++/Go).")

	// 4. Deployment
	fmt.Println("   4. Deployment:")
	fmt.Println("      - Java: Usually requires a JRE and often an App Server (Tomcat/Jetty) for WARs.")
	fmt.Println("      - Go: Produces a native executable; dynamic system-library dependencies may still apply.")
	fmt.Println("      - Binary is cross-compiled easily: GOOS=linux GOARCH=amd64 go build.")
}
