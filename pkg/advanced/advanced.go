// Package advanced provides a deep dive into Go's advanced features, // Package declaration for the advanced features demonstration.
// contrasting them with equivalent concepts in Java. // Explaining the purpose of this package in relation to Java.
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
// Special Agreements & Tooling (Go's unique mechanisms): // Section for Go-specific features.
//
//  1. `//go:generate`: Tool-driven code generation (like Annotation Processors). // Explaining code generation in Go.
//  2. Build Tags (`//go:build`): Conditional compilation. // Explaining conditional compilation in Go.
//     Go has NO C++ style preprocessor (#define). Build tags are used for // Clarifying the lack of a preprocessor.
//     different OS/Arch logic or feature flags at the file level. // Detail on build tags usage.
//  3. `internal` packages: Strong encapsulation (only siblings/parent can import). // Explaining the 'internal' package feature.
//  4. `init()` functions: Automatic per-package setup (like static blocks). // Explaining package initialization functions.
//  5. `iota`: Auto-incrementing counter for constants. Similar to enums, // Explaining the 'iota' keyword.
//     but often used for bitmasks (like C++ #define flags). // Detail on iota usage.
//  6. Struct Tags: Compile-time metadata for reflection (like Jackson annotations). // Explaining struct tags for metadata.
//  7. Records: Go's `struct` is the closest thing to Java `record` (data-focused). // Comparing structs to Java records.
//  8. Bean Agreement: Go prefers direct field access (exported fields) over // Explaining Go's preference for direct access.
//     Getters/Setters unless validation or abstraction is required. // Detail on getter/setter usage in Go.
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
var embeddedFile string // Declaring a variable to hold the content of the embedded file.

// ---------------------------------------------------------------------------
// 8. Struct Tags and JSON (essential for web development) // Section header for Struct Tags and JSON.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Struct tags are like Annotations in Java (e.g., Jackson `@JsonProperty`). // Comparing struct tags to annotations.
// - They are metadata that can be read at runtime via reflection. // Describing how they are used at runtime.
type PersonWithTags struct { // Defining a struct with metadata tags.
	FirstName string `json:"first_name"`    // Field FirstName mapped to "first_name" in JSON.
	LastName  string `json:"last_name"`     // Field LastName mapped to "last_name" in JSON.
	Age       int    `json:"age,omitempty"` // Field Age mapped to "age" and omitted if empty.
	Secret    string `json:"-"`             // Field Secret is ignored during JSON marshaling.
}

/**
 * ===========================================================================
 * Advanced Go Language Features (concepts familiar to Java developers)
 * ===========================================================================
 */ // Block comment for the advanced features section.

// ---------------------------------------------------------------------------
// 1. Embedding (Go's composition – comparable to Java inheritance) // Section header for Embedding.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Go does not have `public`, `private`, or `protected`. // Clarifying Go's access control.
// - "Exporting" is determined by the first letter of the name: // Explaining naming-based visibility.
//   - Uppercase (e.g., `Describe`): Public / Exported (visible outside package). // Upper case means public.
//   - Lowercase (e.g., `secret`): Private / Unexported (visible only in package). // Lower case means private.
//
// - Embedding is composition that looks like inheritance. // Explaining the nature of embedding.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Go doesn't have classes or inheritance. // Reiterating the lack of classes/inheritance.
// - This is a simple struct that will be "embedded" in others. // Describing the purpose of Animal.
type Animal struct {
	Name string // Field Name of type string.
	Legs int    // Field Legs of type int.
}

// Describe returns a string representation of the animal. // Comment for the Describe method.
// Equivalent to a method in a Java base class. // Comparing to Java methods.
func (a Animal) Describe() string { // Method Describe with Animal value receiver.
	return fmt.Sprintf("%s (%d legs)", a.Name, a.Legs)
}

// For a Java developer: // Explanation targeted at Java developers.
//   - This is called "Struct Embedding" or "Anonymous Fields". // Introducing struct embedding terminology.
//   - It's composition, but the methods of the embedded struct are "promoted" // Explaining method promotion.
//   - Java equivalent: `class Dog extends Animal`. // Comparing to Java inheritance.
type Dog struct {
	Animal        // Embedding Animal struct into Dog.
	Breed  string // Additional field Breed of type string.
}

// Speak returns the sound the dog makes. // Comment for the Speak method.
// Equivalent to a method in a subclass in Java. // Comparing to subclass methods.
func (d Dog) Speak() string { // Method Speak with Dog value receiver.
	return fmt.Sprintf("%s says: Woof!", d.Name)
}

type ServiceDog struct {
	Dog           // Embedding Dog struct into ServiceDog.
	CertID string // Additional field CertID for service dog certification.
}

// Describe returns a string representation of the service dog, overriding the embedded Animal's Describe. // Comment for overridden Describe method.
//
// For a Java developer: // Explanation targeted at Java developers.
// - This is similar to `@Override` in Java. // Comparing to Java @Override.
// - If you call `sd.Describe()`, this method is used. // Explaining method resolution.
// - You can still access the "parent" via `sd.Dog.Describe()`. // Explaining parent method access.
// - Java equivalent: `@Override public String describe()` // Comparing to Java equivalent.
func (sd ServiceDog) Describe() string { // Method Describe with ServiceDog value receiver.
	return fmt.Sprintf("%s [Service Dog, Cert: %s]", sd.Dog.Describe(), sd.CertID)
}

// ---------------------------------------------------------------------------
// 2. Enum pattern with iota (comparable to Java enum) // Section header for Enum pattern.
// ---------------------------------------------------------------------------

// Color acts as a type-safe enum using iota. // Comment for the Color type.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Go doesn't have a formal `enum` keyword. // Clarifying the lack of an enum keyword.
// - We use `type MyType int` and a `const` block with `iota`. // Explaining the Go enum pattern.
// - `iota` is a special constant that auto-increments within a const block. // Describing iota's behavior.
// - This is similar to `public static final int` but with type safety. // Comparing to Java static finals.
type Color int // Defining Color as an int type.

const ( // Starting a constant block.
	Red    Color = iota // Red is the first color, iota starts at 0.
	Green               // Green is the next color, iota increments to 1.
	Blue                // Blue is the next color, iota increments to 2.
	Yellow              // Yellow is the next color, iota increments to 3.
)

// String implements the `fmt.Stringer` interface for the Color type. // Comment for the String method.
// Java equivalent: `public String toString()` on an Enum. // Comparing to Java toString.
func (c Color) String() string { // Method String with Color value receiver.
	return [...]string{"Red", "Green", "Blue", "Yellow"}[c]
}

// IsWarm returns true if the color is considered a warm color. // Comment for the IsWarm method.
func (c Color) IsWarm() bool { // Method IsWarm with Color value receiver.
	return c == Red || c == Yellow
}

// Java equivalent: `Color.values()` // Comparing to Java enum values.
func AllColors() []Color { // Function AllColors returning a slice of Color.
	return []Color{Red, Green, Blue, Yellow}
}

// ---------------------------------------------------------------------------
// 3. Functional programming patterns (closures, currying, memoization) // Section header for Functional programming.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `f.compose(g)` on a `Function`. // Comparing to Java function composition.
// - Go functions are first-class citizens and support closures. // Explaining first-class functions in Go.
func Compose[T any](f, g func(T) T) func(T) T { // Generic Compose function with type T.
	return func(x T) T { return f(g(x)) }
}

// For a Java developer: // Explanation targeted at Java developers.
func Curry2[A, B, C any](fn func(A, B) C) func(A) func(B) C { // Generic Curry2 function with types A, B, C.
	return func(a A) func(B) C {
		return func(b B) C {
			return fn(a, b) // Calling the original function with a and b.
		}
	}
}

// For a Java developer: // Explanation targeted at Java developers.
// - Similar to caching results in a `Map` to avoid expensive computations. // Comparing to caching in Java.
func Memoize[K comparable, V any](fn func(K) V) func(K) V { // Generic Memoize function with types K and V.
	cache := make(map[K]V) // Creating a map to store cached results.
	var mu sync.Mutex      // Declaring a mutex for thread-safe access to the cache.
	return func(key K) V {
		mu.Lock() // Locking the mutex to protect the cache.
		defer mu.Unlock()
		if val, ok := cache[key]; ok { // Checking if the result is already in the cache.
			return val
		}
		val := fn(key)   // Calling the original function if not in cache.
		cache[key] = val // Storing the result in the cache.
		return val
	}
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: Chaining multiple `Function.andThen()` calls. // Comparing to Java andThen.
func Pipeline[T any](value T, fns ...func(T) T) T { // Generic Pipeline function with type T.
	for _, fn := range fns { // Iterating over each function in the fns slice.
		value = fn(value) // Applying the function to the current value.
	}
	return value
}

// ---------------------------------------------------------------------------
// 4. Defer / Panic / Recover (comparable to Java try-catch-finally) // Section header for Defer/Panic/Recover.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Go idiom: `panic` is for truly exceptional errors (like `OutOfMemoryError`). // Explaining panic usage.
// - `recover` is used in a `defer` block to catch a panic and handle it gracefully. // Explaining recover usage.
// - Java equivalent: `try-catch` on a `ArithmeticException`. // Comparing to Java try-catch.
func SafeDivide(a, b int) (result int, err error) { // SafeDivide function returning result and error.
	defer func() {
		if r := recover(); r != nil { // Checking if a panic occurred and recovering from it.
			err = fmt.Errorf("recovered panic: %v", r) // Setting the error if a panic was recovered.
		}
	}() // Executing the deferred function at the end of SafeDivide.
	// Integer division by zero panics in Go // Explaining that division by zero causes a panic.
	return a / b, nil
}

// For a Java developer: // Explanation targeted at Java developers.
// - Go idiom: `defer` is the standard way to ensure a resource is closed. // Explaining defer for cleanup.
// - Multiple `defer` calls are executed in Last-In-First-Out (LIFO) order. // Explaining defer execution order.
// - Java equivalent: `try-with-resources` block. // Comparing to Java try-with-resources.
func WithResource(name string) []string { // WithResource function returning a log of operations.
	var log []string                                             // Initializing a slice of strings to store log messages.
	log = append(log, fmt.Sprintf("Opening resource: %s", name)) // Logging the resource opening.
	defer func() {                                               // Deferring a function for cleanup.
		log = append(log, fmt.Sprintf("Closing resource: %s (via defer)", name)) // Logging the resource closing.
	}()
	log = append(log, fmt.Sprintf("Using resource: %s", name)) // Logging the resource usage.
	return log
}

// ---------------------------------------------------------------------------
// 5. Reflection (comparable to Java reflection API) // Section header for Reflection.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: Using `object.getClass()`, `class.getDeclaredFields()`, etc. // Comparing to Java reflection.
// - Go's `reflect` package is powerful but should be used sparingly (it's slow and not type-safe). // Cautioning about reflection usage.
// - This demonstrates `reflect.TypeOf` and `reflect.ValueOf`. // Highlighting the key reflection functions.
func DescribeType(v interface{}) string { // DescribeType function taking an empty interface.
	t := reflect.TypeOf(v)    // Getting the type information of the value.
	val := reflect.ValueOf(v) // Getting the value information of the value.
	var sb strings.Builder    // Initializing a strings.Builder for efficient string concatenation.

	sb.WriteString(fmt.Sprintf("Type: %s, Kind: %s", t.Name(), t.Kind())) // Writing the type name and kind to the builder.

	if t.Kind() == reflect.Struct {
		sb.WriteString(fmt.Sprintf(", Fields: %d", t.NumField()))
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)                                                                 // Getting field metadata.
			fieldVal := val.Field(i)                                                            // Getting the value of the field.
			sb.WriteString(fmt.Sprintf("\n  - %s (%s) = %v", field.Name, field.Type, fieldVal)) // Writing field details to the builder.
		}
	}
	return sb.String()
}

// For a Java developer: // Explanation targeted at Java developers.
// - Similar to how some libraries use reflection to convert POJOs to JSON maps (like Jackson's `objectMapper.convertValue`). // Comparing to POJO-to-map conversion.
// - Demonstrates field visibility (only exported/capitalized fields are processed). // Noting field visibility rules.
func StructToMap(v interface{}) map[string]interface{} { // StructToMap function returning a map.
	result := make(map[string]interface{}) // Initializing a map to hold the struct fields.
	val := reflect.ValueOf(v)
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return result
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)     // Getting field metadata.
		if field.IsExported() { // Checking if the field is exported (public).
			result[field.Name] = val.Field(i).Interface() // Adding the exported field to the result map.
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// 6. Type constraints and type sets (advanced generics) // Section header for Type constraints.
// ---------------------------------------------------------------------------

// Number constrains to all built-in numeric types. // Comment for the Number interface.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Go's type constraints are similar to Java's `extends Number` in generics. // Comparing to Java generic constraints.
// - The `~` symbol allows types that are underlyingly these types (type aliasing). // Explaining the tilde operator in constraints.
type Number interface { // Defining the Number interface as a type constraint.
	~int | ~int8 | ~int16 | ~int32 | ~int64 | // Allowing various signed integer types.
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | // Allowing various unsigned integer types.
		~float32 | ~float64 // Allowing floating-point types.
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `public <T extends Number> T sum(List<T> items)` // Comparing to Java generic sum.
// - Demonstrates type constraints in Go generics. // Highlighting Go generics usage.
func Sum[T Number](items []T) T { // Generic Sum function with Number constraint.
	var total T               // Declaring a variable of type T to hold the total sum.
	for _, v := range items { // Iterating through each item in the slice.
		total += v // Adding the item's value to the total.
	}
	return total
}

// Java equivalent: `Collections.min()` // Comparing to Java Collections.min.
func Min[T Number](items []T) T { // Generic Min function with Number constraint.
	if len(items) == 0 { // Checking if the slice is empty.
		var zero T // Declaring a zero value of type T.
		return zero
	}
	m := items[0]                 // Initializing m with the first element of the slice.
	for _, v := range items[1:] { // Iterating through the rest of the elements.
		if v < m { // Checking if the current element is smaller than m.
			m = v // Updating m if a smaller element is found.
		}
	}
	return m
}

// Java equivalent: `Collections.max()` // Comparing to Java Collections.max.
func Max[T Number](items []T) T { // Generic Max function with Number constraint.
	if len(items) == 0 { // Checking if the slice is empty.
		var zero T // Declaring a zero value of type T.
		return zero
	}
	m := items[0]                 // Initializing m with the first element of the slice.
	for _, v := range items[1:] { // Iterating through the rest of the elements.
		if v > m { // Checking if the current element is larger than m.
			m = v // Updating m if a larger element is found.
		}
	}
	return m
}

// Clamp ensures that a value stays within the range [lo, hi].
// Java equivalent: `Math.clamp()` (available since Java 21)
// Comparing to Java Math.clamp.
func Clamp[T Number](value, lo, hi T) T { // Generic Clamp function with Number constraint.
	if value < lo { // Checking if the value is below the lower bound.
		return lo
	}
	if value > hi { // Checking if the value is above the upper bound.
		return hi
	}
	return value
}

// ---------------------------------------------------------------------------
// 7. Concurrent patterns (advanced – comparable to java.util.concurrent) // Section header for Concurrent patterns.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `ExecutorService.invokeAll()`, parallel streams, or `CompletableFuture.allOf()`. // Comparing to Java concurrency tools.
// - Go idiom: Use a semaphore (buffered channel) to limit concurrency and a `sync.WaitGroup` to wait for completion. // Explaining Go concurrency idioms.
// - This demonstrates how to combine goroutines, channels, and wait groups for complex concurrency patterns. // Noting the combination of techniques.
func FanOut[T any, R any](items []T, workers int, fn func(T) R) []R { // Generic FanOut function with types T and R.
	type indexed struct { // Local struct to track results with their original indices.
		idx    int // Original index of the item.
		result R   // Computed result for the item.
	}
	ch := make(chan indexed, len(items)) // Creating a buffered channel for results.
	sem := make(chan struct{}, workers)  // Creating a semaphore channel to limit concurrent workers.

	var wg sync.WaitGroup        // Declaring a WaitGroup to wait for all goroutines to finish.
	for i, item := range items { // Iterating through each item in the input slice.
		wg.Add(1)                 // Incrementing the WaitGroup counter.
		go func(idx int, val T) { // Launching a goroutine for each item.
			defer wg.Done()            // Ensuring Done is called when the goroutine finishes.
			sem <- struct{}{}          // Acquiring a slot in the semaphore.
			result := fn(val)          // Executing the provided function on the item.
			<-sem                      // Releasing the slot in the semaphore.
			ch <- indexed{idx, result} // Sending the result along with its index to the results channel.
		}(i, item) // Passing loop variables to the goroutine closure.
	}
	go func() { // Launching a separate goroutine to close the results channel.
		wg.Wait() // Waiting for all worker goroutines to complete.
		close(ch)
	}() // Executing the closer goroutine.

	results := make([]R, len(items)) // Initializing the final results slice.
	for r := range ch {              // Reading from the results channel until it is closed.
		results[r.idx] = r.result // Placing the result at its original index.
	}
	return results
}

// ---------------------------------------------------------------------------
// 11. Functional Options Pattern (Common Go pattern for configuration) // Section header for Functional Options Pattern.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - This is an alternative to the Builder pattern. // Comparing to the Java Builder pattern.
// - It is very common in Go libraries for configuring clients or servers. // Noting its prevalence in Go.
type DatabaseConfig struct {
	Host     string        // Host field for database connection string.
	Port     int           // Port field for database connection.
	Timeout  time.Duration // Timeout duration for database operations.
	MaxConns int           // Maximum number of connections in the pool.
}

// DBOption is a function type that modifies DatabaseConfig. // Comment for the DBOption type.
type DBOption func(*DatabaseConfig) // Defining DBOption as a function taking a pointer to DatabaseConfig.

// Java equivalent: `builder.setHost(host)` // Comparing to Java builder method.
func WithHost(host string) DBOption { // Function WithHost returning a DBOption closure.
	return func(c *DatabaseConfig) {
		c.Host = host // Setting the host field in the config.
	}
}

// Java equivalent: `builder.setPort(port)` // Comparing to Java builder method.
func WithPort(port int) DBOption { // Function WithPort returning a DBOption closure.
	return func(c *DatabaseConfig) {
		c.Port = port // Setting the port field in the config.
	}
}

// Java equivalent: `builder.setTimeout(t)` // Comparing to Java builder method.
func WithTimeout(t time.Duration) DBOption { // Function WithTimeout returning a DBOption closure.
	return func(c *DatabaseConfig) {
		c.Timeout = t // Setting the timeout field in the config.
	}
}

// Java equivalent: `new DatabaseConnector(host, port, timeout, maxConns)` // Comparing to Java constructor.
func NewDatabaseConnector(opts ...DBOption) *DatabaseConfig { // Function NewDatabaseConnector with variadic options.
	// Default values // Setting up default configuration values.
	config := &DatabaseConfig{ // Initializing the config with defaults.
		Host:     "localhost",      // Default host is localhost.
		Port:     5432,             // Default port is 5432 (Postgres).
		Timeout:  30 * time.Second, // Default timeout is 30 seconds.
		MaxConns: 10,               // Default max connections is 10.
	}

	for _, opt := range opts { // Iterating through each provided option.
		opt(config) // Applying the option function to the config.
	}
	return config
}

// ---------------------------------------------------------------------------
// 12. Atomic Operations (Low-level synchronization) // Section header for Atomic Operations.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - This is equivalent to `AtomicLong` or `AtomicInteger`. // Comparing to Java atomic types.
// - It is more performant than using a `sync.Mutex` for simple counter operations. // Noting performance benefits over mutexes.
type SafeCounter struct {
	value atomic.Int64 // Using atomic.Int64 for thread-safe value storage.
}

// Increment safely increments the counter. // Comment for the Increment method.
// Java equivalent: `atomicLong.incrementAndGet()` // Comparing to Java atomic increment.
func (c *SafeCounter) Increment() { // Method Increment with SafeCounter pointer receiver.
	c.value.Add(1) // Atomically adding 1 to the counter value.
}

// Value returns the current counter value. // Comment for the Value method.
// Java equivalent: `atomicLong.get()` // Comparing to Java atomic get.
func (c *SafeCounter) Value() int64 { // Method Value with SafeCounter pointer receiver.
	return c.value.Load() // Atomically loading the current counter value.
}

// ---------------------------------------------------------------------------
// 13. Compile-time Interface Satisfaction Check // Section header for Compile-time Interface Check.
// ---------------------------------------------------------------------------

// Logger is a simple logging interface. // Comment for the Logger interface.
type Logger interface { // Defining the Logger interface.
	Log(message string) // Method Log taking a string message.
}

// ConsoleLogger implements the Logger interface. // Redundant comment for ConsoleLogger.
type ConsoleLogger struct{} // Defining the ConsoleLogger struct (empty struct).

// Log prints the message to the console. // Comment for the Log method.
// Equivalent to `System.out.println()` in Java. // Comparing to Java's System.out.println.
func (c ConsoleLogger) Log(message string) { // Method Log with ConsoleLogger value receiver.
	fmt.Println("   [Logger]: " + message) // Printing the message with a prefix to the console.
}

// This line ensures that ConsoleLogger implements Logger at compile time. // Comment explaining the interface check.
// If it doesn't, the compiler will throw an error. // Noting the compiler error benefit.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: None (Java checks at compile time based on `implements` keyword). // Contrasting with Java's explicit implements.
// - Go idiom: Use a blank identifier (`_`) and a type conversion to ensure an interface is satisfied. // Explaining the Go idiom for interface check.
var _ Logger = (*ConsoleLogger)(nil) // Ensuring ConsoleLogger satisfies Logger interface at compile time.

// ---------------------------------------------------------------------------
// 14. Generic Concurrent Cache (RWMutex + Generics) // Section header for Generic Concurrent Cache.
// ---------------------------------------------------------------------------

// For a Java developer: // Explanation targeted at Java developers.
// - Similar to `ConcurrentHashMap<K, V>`. // Comparing to Java's ConcurrentHashMap.
// - `sync.RWMutex` allows multiple concurrent readers but only one writer. // Explaining the RWMutex behavior.
// - Unlike Java's `ConcurrentHashMap`, this is an "explicit" lock-based implementation. // Noting the explicit locking nature.
type ConcurrentCache[K comparable, V any] struct {
	mu    sync.RWMutex // Mutex for protecting map access.
	items map[K]V      // Underlying map to store cache items.
}

// Java equivalent: `new ConcurrentHashMap<>()` // Comparing to Java ConcurrentHashMap instantiation.
func NewConcurrentCache[K comparable, V any]() *ConcurrentCache[K, V] { // Generic function returning a new cache.
	return &ConcurrentCache[K, V]{ // Initializing and returning the cache pointer.
		items: make(map[K]V), // Creating the underlying map.
	}
}

// Set adds or updates a value in the cache. // Comment for the Set method.
// Java equivalent: `cache.put(key, value)` // Comparing to Java map put.
func (c *ConcurrentCache[K, V]) Set(key K, value V) { // Method Set with ConcurrentCache pointer receiver.
	c.mu.Lock()          // Acquiring a write lock to update the map.
	defer c.mu.Unlock()  // Deferring the release of the write lock.
	c.items[key] = value // Adding or updating the value in the map.
}

// Get retrieves a value from the cache. // Comment for the Get method.
// Java equivalent: `cache.get(key)` // Comparing to Java map get.
func (c *ConcurrentCache[K, V]) Get(key K) (V, bool) { // Method Get with ConcurrentCache pointer receiver.
	c.mu.RLock()            // Acquiring a read lock for concurrent reading.
	defer c.mu.RUnlock()    // Deferring the release of the read lock.
	val, ok := c.items[key] // Retrieving the value and existence status from the map.
	return val, ok
}

// ---------------------------------------------------------------------------
// 15. Context-Aware Multi-stage Pipeline (Advanced Concurrency) // Section header for Advanced Concurrency Pipeline.
// ---------------------------------------------------------------------------

type pipelineResult struct { // Defining a struct to hold stage results and errors.
	val int   // Integer value result from a pipeline stage.
	err error // Error encountered during the stage, if any.
}

// For a Java developer: // Explanation targeted at Java developers.
// - Java equivalent: `CompletableFuture` chain or a reactive stream (RxJava/Project Reactor) with cancellation support. // Comparing to Java reactive/future chains.
// - Demonstrates multi-stage processing using channels and context for graceful shutdown. // Noting the use of channels and context.
func AdvancedPipeline(ctx context.Context, nums []int) (int, error) { // AdvancedPipeline function with context.
	if err := ctx.Err(); err != nil { // Checking if the context is already cancelled.
		return 0, err
	}

	// Stage 1: Square the numbers // Comment for the first stage.
	stage1 := make(chan pipelineResult) // Creating a channel for the first stage results.
	go func() {                         // Launching a goroutine for the first stage.
		defer close(stage1)      // Ensuring the stage channel is closed when the goroutine ends.
		for _, n := range nums { // Iterating through the input numbers.
			select { // Using select to handle both processing and cancellation.
			case <-ctx.Done(): // Case for context cancellation.
				return // Exiting the goroutine if cancelled.
			case stage1 <- pipelineResult{val: n * n}: // Case for sending the result to the channel.
			}
		}
	}() // Executing the stage 1 goroutine.

	// Stage 2: Filter even squares and sum them // Comment for the second stage.
	sum := 0                  // Initializing the sum variable.
	for res := range stage1 { // Reading results from the first stage channel.
		if res.err != nil { // Checking if an error occurred in the previous stage.
			return 0, res.err
		}
		if res.val%2 == 0 { // Filtering for even squared values.
			sum += res.val // Adding the even value to the running sum.
		}

		// Check for context cancellation in the consumer loop too // Explaining the need for cancellation check in consumer.
		select { // Using select for non-blocking cancellation check.
		case <-ctx.Done(): // Case for context cancellation.
			return 0, ctx.Err()
		default: // Default case to continue processing if not cancelled.
		}
	}

	return sum, nil
}

// ---------------------------------------------------------------------------
// 16. Special Agreements (Tooling) // Section header for Special Agreements and Tooling.
// ---------------------------------------------------------------------------

//go:generate echo "Running custom code generation tool..."
// ^ Comment for the go:generate directive above.

// For a Java developer: // Explanation targeted at Java developers.
// - Go doesn't use the JavaBeans pattern (get/set everywhere). // Noting the difference from JavaBeans.
// - We only use getters/setters if there's actual logic (like validation). // Explaining when to use getters/setters in Go.
// - Exported fields (starts with uppercase) are the idiomatic way to expose data. // Defining the Go way to expose data.
type Account struct {
	id      int    // unexported (private to package) field id.
	balance int    // unexported (private to package) field balance.
	Owner   string // exported (public) field Owner.
}

// GetBalance is a Getter. It uses a value receiver. // Comment for the GetBalance method.
// Java equivalent: `public int getBalance()` // Comparing to Java getter.
func (a Account) GetBalance() int { // Method GetBalance with Account value receiver.
	return a.balance
}

// Deposit is a Setter. It uses a pointer receiver to modify the balance. // Comment for the Deposit method.
// Java equivalent: `public void deposit(int amount)` // Comparing to Java setter.
func (a *Account) Deposit(amount int) { // Method Deposit with Account pointer receiver.
	if amount > 0 { // Validating that the deposit amount is positive.
		a.balance += amount // Incrementing the balance by the deposit amount.
	}
}

// For a Java developer: // Explanation targeted at Java developers.
// - Similar to a Java 14+ `record`. // Comparing to Java records.
type UserRecord struct {
	Username string // Field Username of type string.
	Email    string // Field Email of type string.
}

// For a Java developer: // Explanation targeted at Java developers.
//   - Go does not have a preprocessor, but it uses "Build Tags" and "Generate" // Explaining alternatives to preprocessors in Go.
//     to achieve similar results. // Describing how similar results are achieved.
func RunAdvancedDemo() { // RunAdvancedDemo function for showcasing features.
	fmt.Println("--- Advanced Language Features Demo ---") // Printing the demo start header.

	// 1. Embedding // Comment for the embedding demonstration.
	fmt.Println("1. Embedding (composition ≈ Java inheritance):")              // Printing the section title.
	dog := Dog{Animal: Animal{Name: "Rex", Legs: 4}, Breed: "German Shepherd"} // Initializing a Dog struct with embedded Animal.
	fmt.Printf("   %s (Breed: %s)\n", dog.Describe(), dog.Breed)               // Printing the dog's description and breed.
	fmt.Printf("   %s\n", dog.Speak())                                         // Printing the dog's vocalization.
	sd := ServiceDog{Dog: dog, CertID: "SD-42"}                                // Initializing a ServiceDog struct embedding Dog.
	fmt.Printf("   %s\n", sd.Describe())                                       // Printing the service dog's overridden description.

	// 2. Enums with iota // Comment for the enum demonstration.
	fmt.Println("2. Enum pattern (iota ≈ Java enum):") // Printing the section title.
	for _, c := range AllColors() {                    // Iterating through all defined colors.
		fmt.Printf("   %s (ordinal=%d, warm=%v)\n", c, c, c.IsWarm()) // Printing color name, ordinal value, and warmth.
	}

	// 3. Functional patterns // Comment for the functional patterns demonstration.
	fmt.Println("3. Functional programming:")                    // Printing the section title.
	double := func(x int) int { return x * 2 }                   // Defining a local anonymous function to double an int.
	addThree := func(x int) int { return x + 3 }                 // Defining a local anonymous function to add three to an int.
	doubleThenAdd := Compose(addThree, double)                   // Composing the two functions using the Compose utility.
	fmt.Printf("   Compose(+3, *2)(5) = %d\n", doubleThenAdd(5)) // Printing the result of function composition.

	curriedAdd := Curry2(func(a, b int) int { return a + b })
	add10 := curriedAdd(10)
	fmt.Printf("   Curry add(10)(5) = %d\n", add10(5)) // Printing the final result of the curried call.

	calls := 0 // Counter to track the number of function calls.
	expensive := Memoize(func(n int) int {
		calls++ // Incrementing the call counter.
		return n * n
	})
	expensive(4)                                                                     // Calling the expensive function for the first time.
	expensive(4)                                                                     // Calling it a second time (should be cached).
	expensive(4)                                                                     // Calling it a third time (should be cached).
	fmt.Printf("   Memoize: square(4) called 3 times, computed %d time(s)\n", calls) // Printing the number of actual computations.

	result := Pipeline("  Hello, World!  ", // Passing a string through a processing pipeline.
		strings.TrimSpace, // Trimming whitespace from both ends.
		strings.ToUpper,   // Converting the string to uppercase.
		func(s string) string { return ">>>" + s + "<<<" }, // Adding custom formatting.
	)
	fmt.Printf("   Pipeline: %s\n", result) // Printing the final pipeline result.

	// 4. Defer / Panic / Recover // Comment for the error handling demonstration.
	fmt.Println("4. Defer/Panic/Recover (≈ try-catch-finally):")  // Printing the section title.
	val, err := SafeDivide(10, 3)                                 // Performing a safe division that handles panics.
	fmt.Printf("   SafeDivide(10, 3) = %d, err = %v\n", val, err) // Printing result and error for valid division.
	val, err = SafeDivide(10, 0)                                  // Performing safe division by zero.
	fmt.Printf("   SafeDivide(10, 0) = %d, err = %v\n", val, err) // Printing result and error for invalid division.
	// Note: WithResource's defer runs before the function returns, // Comment explaining defer timing.
	// so the "Closing" line is appended before we see the log. // Noting log message order.
	// We show the concept differently: // Explaining the alternative demonstration approach.
	fmt.Printf("   Defer ordering: ") // Printing a label for defer ordering.
	func() {                          // Anonymous function to demonstrate LIFO defer execution.
		defer fmt.Print("third ")  // Third defer should execute last of these three? No, LIFO.
		defer fmt.Print("second ") // Second defer.
		defer fmt.Print("first ")  // First defer to execute.
	}()
	fmt.Println("(defers execute LIFO)") // Printing a note about LIFO order.

	// 5. Reflection // Comment for the reflection demonstration.
	fmt.Println("5. Reflection (≈ Java reflection API):") // Printing the section title.
	type User struct {                                    // Defining a local User struct for reflection demo.
		Name  string // Name field.
		Age   int    // Age field.
		Email string // Email field.
	}
	u := User{Name: "Alice", Age: 30, Email: "alice@example.com"} // Initializing a User instance.
	fmt.Printf("   %s\n", DescribeType(u))                        // Printing a reflection-based description of the User instance.
	m := StructToMap(u)                                           // Converting the struct to a map using reflection.
	fmt.Printf("   StructToMap: %v\n", m)                         // Printing the resulting map.

	// 6. Type constraints (advanced generics) // Comment for the generics demonstration.
	fmt.Println("6. Type constraints (Number interface):")      // Printing the section title.
	ints := []int{3, 1, 4, 1, 5, 9, 2, 6}                       // Initializing a slice of integers.
	fmt.Printf("   Sum(%v) = %d\n", ints, Sum(ints))            // Printing the sum of integers.
	fmt.Printf("   Min = %d, Max = %d\n", Min(ints), Max(ints)) // Printing the min and max of integers.
	fmt.Printf("   Clamp(15, 0, 10) = %d\n", Clamp(15, 0, 10))  // Printing the result of clamping an integer.
	floats := []float64{3.14, 2.71, 1.41}                       // Initializing a slice of floats.
	fmt.Printf("   Sum(%v) = %.2f\n", floats, Sum(floats))      // Printing the sum of floats.

	// 7. Concurrent fan-out // Comment for the concurrency demonstration.
	fmt.Println("7. Concurrent fan-out (≈ ExecutorService):")      // Printing the section title.
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8}                        // Initializing input data.
	squares := FanOut(inputs, 3, func(n int) int { return n * n }) // Processing inputs concurrently with 3 workers.
	fmt.Printf("   FanOut squares: %v\n", squares)                 // Printing the final results.

	// 8. Struct Tags and JSON // Comment for the struct tags demonstration.
	fmt.Println("8. Struct Tags and JSON (common in REST APIs):")                            // Printing the section title.
	p := PersonWithTags{FirstName: "John", LastName: "Doe", Age: 30, Secret: "password123"}  // Initializing a struct with tags.
	data, _ := json.MarshalIndent(p, "   ", "  ")                                            // Marshaling the struct to JSON with indentation.
	fmt.Printf("   Marshaled: %s\n", string(data))                                           // Printing the marshaled JSON string.
	fmt.Println("   (Note: 'Secret' is ignored and 'Age' would be omitted if it were zero)") // Noting the effect of struct tags.

	// 9. Embed (Go 1.16+) // Comment for the embed demonstration.
	fmt.Println("9. Embed package (embedding static assets):")   // Printing the section title.
	fmt.Printf("   Embedded file content:\n   %s", embeddedFile) // Printing the content of the embedded file.

	// 10. Advanced Iota (Bitmask/Flags) // Comment for the advanced iota demonstration.
	fmt.Println("10. Advanced Iota (Bitmask flags):") // Printing the section title.
	type Permission int                               // Defining Permission as an int type.
	const (                                           // Starting a constant block for permissions.
		Read    Permission = 1 << iota // 1 << 0 = 1.
		Write                          // 1 << 1 = 2.
		Execute                        // 1 << 2 = 4.
	)
	pReadWrite := Read | Write                                         // Creating a combined permission using bitwise OR.
	fmt.Printf("   Permissions: %d (Read:%v, Write:%v, Execute:%v)\n", // Printing permission details.
		pReadWrite, pReadWrite&Read != 0, pReadWrite&Write != 0, pReadWrite&Execute != 0) // Checking individual bits.

	// 11. Functional Options Pattern // Comment for the functional options demonstration.
	fmt.Println("11. Functional Options Pattern (Go configuration):") // Printing the section title.
	dbConfig := NewDatabaseConnector(                                 // Configuring a database connection using functional options.
		WithHost("db.example.com"),  // Setting the host.
		WithPort(5433),              // Setting the port.
		WithTimeout(10*time.Second), // Setting the timeout.
	)
	fmt.Printf("   DBConfig: %+v\n", dbConfig) // Printing the final config object.

	// 12. Atomic operations (lock-free synchronization) // Comment for the atomic demonstration.
	fmt.Println("12. Atomic operations (lock-free synchronization):") // Printing the section title.
	counter := &SafeCounter{}                                         // Initializing a safe counter.
	var wg sync.WaitGroup                                             // Declaring a WaitGroup to wait for all goroutines.
	for i := 0; i < 100; i++ {                                        // Launching 100 concurrent increments.
		wg.Add(1)   // Incrementing WaitGroup counter.
		go func() { // Launching a goroutine.
			defer wg.Done()     // ensuring Done is called.
			counter.Increment() // Safely incrementing the counter.
		}()
	}
	wg.Wait()                                                                           // Waiting for all increments to finish.
	fmt.Printf("   SafeCounter value (100 parallel increments): %d\n", counter.Value()) // Printing the final counter value.

	// 13. Compile-time Interface Check // Comment for the interface check demonstration.
	fmt.Println("13. Compile-time Interface Check:")         // Printing the section title.
	var _ Logger = (*ConsoleLogger)(nil)                     // Validated at compile time to ensure ConsoleLogger satisfies Logger.
	logger := &ConsoleLogger{}                               // Initializing a ConsoleLogger.
	logger.Log("Hello from the interface-validated logger!") // Using the logger.

	// 14. Generic Concurrent Cache (RWMutex + Generics) // Comment for the generic cache demonstration.
	fmt.Println("14. Generic Concurrent Cache (RWMutex + Generics):") // Printing the section title.
	cache := NewConcurrentCache[string, int]()                        // Initializing a new concurrent cache.
	cache.Set("Go", 2009)                                             // Storing a value in the cache.
	cache.Set("Java", 1995)                                           // Storing another value.
	if val, ok := cache.Get("Go"); ok {                               // Retrieving a value from the cache.
		fmt.Printf("   Cache: Go was released in %d\n", val) // Printing the retrieved value.
	}

	// 15. Advanced Context-Aware Pipeline // Comment for the advanced pipeline demonstration.
	fmt.Println("15. Context-Aware Multi-stage Pipeline:")                         // Printing the section title.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond) // Creating a context with a timeout.
	defer cancel()                                                                 // Ensuring the context is cancelled to release resources.
	pipelineSum, err := AdvancedPipeline(ctx, []int{1, 2, 3, 4, 5, 6})             // Executing the advanced pipeline.
	if err != nil {                                                                // Checking for errors during pipeline execution.
		fmt.Printf("   Pipeline error: %v\n", err) // Printing any encountered error.
	} else { // If no error occurred.
		fmt.Printf("   Pipeline Sum of even squares: %d\n", pipelineSum) // Printing the calculated sum.
	}

	// 16. Java Developer FAQ // Comment for the FAQ section.
	fmt.Println("16. Java Developer FAQ (Getters/Setters, Records, Annotations):") // Printing the section title.
	RunJavaDeveloperFAQ()                                                          // Running the Java developer FAQ demonstration.

	// 17. Build System Comparison // Comment for the build system comparison section.
	fmt.Println("17. Build System Comparison (Go vs Maven):") // Printing the section title.
	RunBuildSystemDemo()                                      // Running the build system comparison demonstration.

	fmt.Println("--- Advanced Language Features Demo End ---") // Printing the demo end footer.
}

// RunJavaDeveloperFAQ provides explicit examples for common Java vs Go questions.
func RunJavaDeveloperFAQ() {
	// a. Getters and Setters (Bean Agreement)
	// Java-ism: Writing GetX() and SetX() for every field.
	// Go way: Use direct field access (Uppercase = Public).
	// Only use methods if you need to encapsulate internal logic or provide a read-only view.
	acc := Account{id: 1, balance: 100, Owner: "Alice"}
	acc.Deposit(50) // Method used for logic (adding to balance).
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
	var buf strings.Builder // Ready to use!
	buf.WriteString("   - Zero-value: Ready to use without a constructor!")
	fmt.Println(buf.String())

	// e. Pointers vs Values
	// Java-ism: Passing everything as a pointer (*T) because Java objects are references.
	// Go way: Pass by value (T) by default. Use pointers ONLY if you need to:
	// 1. Mutate the original object.
	// 2. Avoid copying a VERY large struct (usually > 64 bytes).
	// 3. Represent 'nil' for an optional value.
	fmt.Println("   - Pointers: Use values by default; use pointers only when mutation or size matters.")

	// f. Records vs Structs
	// Java 14+ Records are immutable by default. Go Structs are mutable unless you use a value receiver.
	u := UserRecord{"gopher", "go@golang.org"}
	fmt.Printf("   - Record-like Struct: %+v\n", u)

	// g. Annotations vs Struct Tags
	// Java Annotations can be processed at compile-time or runtime.
	// Struct tags are only available via reflection at runtime.
	fmt.Println("   - Annotations: Use Struct Tags for metadata (Runtime) and go:generate (Compile-time).")

	// h. Threads vs Goroutines
	// Goroutines are much cheaper. You can run 100k+ goroutines on a laptop.
	fmt.Println("   - Threads: Java Threads are heavy (~1MB stack); Go Goroutines are light (~2KB).")

	// i. No Magic (AOP, Annotations, Dependency Injection)
	// Java-ism: Relying on reflection/proxies for business logic (e.g. @Transactional).
	// Go way: Be explicit. Pass dependencies as arguments; use closures or decorators for middleware.
	// Go avoids "magic" that happens behind the scenes.
	fmt.Println("   - No Magic: Go prefers explicit code over heavy AOP or reflection-based frameworks.")
}

func RunBuildSystemDemo() {
	// 1. Build Tool // Comment for the build tool section.
	fmt.Println("   1. Build Tool: Go uses the 'go' tool (CLI-first).")                                                      // Stating the primary build tool.
	fmt.Println("      - 'go build': Compiles the project (≈ 'mvn package'). Produces a self-contained binary.")             // Explaining go build.
	fmt.Println("      - 'go test': Runs all tests (≈ 'mvn test'). Native support for unit tests, benchmarks, and fuzzing.") // Explaining go test.
	fmt.Println("      - 'go run': Compiles and runs (≈ 'mvn exec:java'). No need for manual compilation during dev.")       // Explaining go run.

	// 2. Dependency Management // Comment for the dependency management section.
	fmt.Println("   2. Dependency Management: Go uses Go Modules ('go.mod').")                                          // Stating the dependency management system.
	fmt.Println("      - No XML! It's a simple text file (≈ 'pom.xml').")                                               // Highlighting the simplicity over Maven's XML.
	fmt.Println("      - 'go mod tidy': Syncs dependencies (≈ Maven Import/Refresh). It also cleans up unused ones.")   // Explaining go mod tidy.
	fmt.Println("      - Dependencies are downloaded to $GOPATH/pkg/mod, shared across projects (≈ ~/.m2/repository).") // Explaining dependency storage.

	// 3. Project Structure and Plugins // Comment for the plugins section.
	fmt.Println("   3. Plugins & Tooling:")                                                         // Stating the plugins and tooling section.
	fmt.Println("      - In Maven, you use plugins (maven-compiler-plugin, etc.).")                 // Describing Maven's plugin model.
	fmt.Println("      - In Go, the 'go' tool is batteries-included (fmt, vet, test, doc, build).") // Highlighting Go's built-in tools.
	fmt.Println("      - For custom tasks, Go devs use a 'Makefile' (common in C/C++/Go).")         // Mentioning Makefiles for custom orchestration.

	// 4. Deployment // Comment for the deployment section.
	fmt.Println("   4. Deployment:")                                                                              // Stating the deployment section.
	fmt.Println("      - Java: Usually requires a JRE and often an App Server (Tomcat/Jetty) for WARs.")          // Describing Java deployment requirements.
	fmt.Println("      - Go: Produces a single, statically-linked binary (≈ a Fat JAR with no JVM requirement).") // Highlighting Go's self-contained binary.
	fmt.Println("      - Binary is cross-compiled easily: GOOS=linux GOARCH=amd64 go build.")                     // Mentioning cross-compilation simplicity.
}
