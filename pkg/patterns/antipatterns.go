// Package patterns showcases GoF design patterns and Go-specific antipatterns.
//
// This file focuses on Go Antipatterns commonly committed by Java developers.
// Coming from the JVM, Java developers carry object-oriented assumptions
// (everything is a heap reference, null checks, thread pools, exception hierarchies)
// that do not map 1:1 to Go's value-oriented runtime, structural interfaces,
// and goroutine concurrency model.
package patterns

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"sync"
	"time"
)

// ===========================================================================
// Antipattern 1: The Nil Interface Trap (Typed Nil Pointer in Interface)
// ===========================================================================
//
// For a Java developer:
// In Java, if a reference variable `ex` is `null`, returning `ex` as an `Exception`
// or `Object` means the caller sees `null`. Checking `if (ex != null)` correctly
// branches to the non-error path.
//
// Go Trap:
// In Go, an interface value is an internal two-word struct: (concreteType, concreteValue).
// An interface equals `nil` ONLY when BOTH its type and its value are `nil`.
// If a concrete pointer `var customErr *CustomValidationError = nil` is returned
// as an `error` interface, the returned interface is `(*CustomValidationError, nil)`.
// Because the type is NOT nil, `err != nil` evaluates to TRUE!
// The caller believes an error occurred and enters the error handling branch.
//
// Idiomatic Go:
// Always return the untyped literal `nil` when there is no error: `return nil`.
// Or declare the variable directly as the interface type: `var err error`.

// CustomValidationError is a concrete domain error.
type CustomValidationError struct {
	Field string
	Msg   string
}

func (e *CustomValidationError) Error() string {
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Msg)
}

// ValidateUsernameAntipattern returns a concrete pointer stored in an error interface.
// Even if isValid is true, the returned error interface is NOT nil!
func ValidateUsernameAntipattern(name string) error {
	var validationErr *CustomValidationError = nil
	if len(name) < 3 {
		validationErr = &CustomValidationError{Field: "username", Msg: "too short"}
	}
	// BUG: Returning a typed nil pointer as an error interface creates an interface
	// value with Type=*CustomValidationError and Value=nil. (err != nil) is TRUE!
	return validationErr
}

// ValidateUsernameIdiomatic explicitly returns literal nil on success.
func ValidateUsernameIdiomatic(name string) error {
	if len(name) < 3 {
		return &CustomValidationError{Field: "username", Msg: "too short"}
	}
	// Idiomatic Go: Return untyped literal nil so the interface itself is nil.
	return nil
}

// ===========================================================================
// Antipattern 2: Passing Mutex by Value (Lock State Copying)
// ===========================================================================
//
// For a Java developer:
// In Java, `synchronized (this)` or `ReentrantLock` are objects allocated on the JVM heap.
// Passing an object reference passes a copy of the pointer; all callers share the
// exact same underlying lock monitor.
//
// Go Trap:
// Structs in Go are value types. If a struct containing a `sync.Mutex` is passed
// by value, or if its methods use a value receiver `func (c Counter) Inc()`, Go makes
// a complete byte-for-byte copy of the struct—including the internal lock state!
// Each caller acquires its own independent copy of the mutex, completely eliminating
// mutual exclusion and resulting in silent data races.
//
// Idiomatic Go:
// Types holding synchronizers (`sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`) MUST
// always be passed by pointer (`*Counter`) and must only use pointer receivers.
// Note: Go's `go vet` tool has a built-in `copylocks` analyzer that flags any method
// receiving or returning `sync.Mutex` by value at compile time!

// SimulatedLock represents a lock state for demonstration.
// (We use Acquire/Release so this educational antipattern demo does not trip `go vet`
// during CI / build checks, which strictly rejects any `sync.Mutex` passed by value).
type SimulatedLock struct {
	heldBy int
}

func (l *SimulatedLock) Acquire(id int) { l.heldBy = id }
func (l *SimulatedLock) Release()       { l.heldBy = 0 }

// ValueCounterAntipattern defines a counter with a value receiver (antipattern).
// Because the receiver `c` is a copy, caller and callee have completely separate locks!
type ValueCounterAntipattern struct {
	lock  SimulatedLock
	count int
}

// Inc uses a value receiver: it copies the struct and its lock on every call!
func (c ValueCounterAntipattern) Inc(id int) ValueCounterAntipattern {
	c.lock.Acquire(id)
	defer c.lock.Release()
	c.count++
	return c
}

// PointerCounterIdiomatic uses pointer receivers to guarantee a single shared lock.
type PointerCounterIdiomatic struct {
	mu    sync.Mutex
	count int
}

// Inc uses a pointer receiver: the caller's mutex is actually locked and count incremented.
func (c *PointerCounterIdiomatic) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Value returns the current count under read synchronization.
func (c *PointerCounterIdiomatic) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ===========================================================================
// Antipattern 3: Goroutine Leaks (Unbuffered Channels & No Cancellation)
// ===========================================================================
//
// For a Java developer:
// In Java, daemon threads or `ExecutorService.shutdown()` can terminate workers.
// Go goroutines have no external cancel method (`thread.stop()` doesn't exist).
//
// Go Trap:
// If a goroutine attempts to send to an unbuffered channel (`ch <- val`), but the
// receiving function returns early (e.g. timeout, first match, or error), the
// goroutine blocks indefinitely. It is never garbage collected, leaking its
// stack, memory, and any system resources it holds.
//
// Idiomatic Go:
// 1. For single-response worker goroutines, use a buffered channel of capacity 1:
//    `make(chan string, 1)`. Even if the receiver abandons the read, the send succeeds.
// 2. For multi-step tasks, accept a `context.Context` and select on `ctx.Done()`.

// QueryDataAntipattern leaks a goroutine if the function returns early.
func QueryDataAntipattern(slow bool) (string, error) {
	ch := make(chan string) // Unbuffered channel

	go func() {
		if slow {
			time.Sleep(50 * time.Millisecond)
		}
		// If caller timed out and returned, this send blocks FOREVER -> Goroutine Leak!
		ch <- "result-from-slow-query"
	}()

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(10 * time.Millisecond):
		return "", errors.New("query timed out")
	}
}

// QueryDataIdiomatic uses a buffered channel so the background goroutine can always exit.
func QueryDataIdiomatic(ctx context.Context, slow bool) (string, error) {
	// Buffered channel of capacity 1 prevents the goroutine from blocking forever
	ch := make(chan string, 1)

	go func() {
		if slow {
			select {
			case <-time.After(50 * time.Millisecond):
			case <-ctx.Done():
				return // Abort early if context cancelled
			}
		}
		ch <- "result-from-slow-query"
	}()

	select {
	case res := <-ch:
		return res, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// ===========================================================================
// Antipattern 4: Interface Pollution (Producer-side / Pre-emptive Interfaces)
// ===========================================================================
//
// For a Java developer:
// In Spring/Java EE, standard architecture mandates:
//   public interface CustomerService { ... }
//   public class CustomerServiceImpl implements CustomerService { ... }
// even when there is only ever one concrete implementation.
//
// Go Trap:
// In Go, interfaces are satisfied implicitly (duck typing). Creating 10-method
// interfaces in the same package as the implementation before any consumer asks
// for it creates rigid coupling, forces unnecessary mock boilerplate, and defeats
// Go's structural typing.
//
// Idiomatic Go:
// "Accept interfaces, return structs."
// 1. The producer package exports concrete structs: `func NewCustomerRepository() *CustomerRepository`.
// 2. The consumer package defines small, focused interfaces (often 1 or 2 methods)
//    declaring only the methods it actually needs.

// Customer is a simple entity.
type Customer struct {
	ID   string
	Name string
}

// CustomerStore is a concrete repository.
type CustomerStore struct {
	data map[string]*Customer
}

// NewCustomerStore returns a concrete pointer, NOT an interface.
func NewCustomerStore() *CustomerStore {
	return &CustomerStore{
		data: map[string]*Customer{
			"c-1": {ID: "c-1", Name: "Acme Corp"},
			"c-2": {ID: "c-2", Name: "Globex Corp"},
		},
	}
}

// GetCustomer retrieves a customer by ID.
func (cs *CustomerStore) GetCustomer(id string) (*Customer, error) {
	c, ok := cs.data[id]
	if !ok {
		return nil, errors.New("customer not found")
	}
	return c, nil
}

// CustomerFinder is a small, consumer-defined interface containing only the method needed.
// Notice it is defined by the caller, not dictated by the CustomerStore package!
type CustomerFinder interface {
	GetCustomer(id string) (*Customer, error)
}

// LookupCustomerName accepts the narrow interface, allowing any mock or alternative implementation.
// It propagates lookup failures rather than turning them into a success-shaped fallback.
func LookupCustomerName(finder CustomerFinder, id string) (string, error) {
	c, err := finder.GetCustomer(id)
	if err != nil {
		return "", fmt.Errorf("lookup customer %q: %w", id, err)
	}
	return c.Name, nil
}

// ===========================================================================
// Antipattern 5: Pointers to Reference Types (*[]T, *map[K]V, *chan T)
// ===========================================================================
//
// For a Java developer:
// Java developers are accustomed to thinking: "Objects are references, but if I pass
// a big collection to a method, I should pass a reference to avoid copying it."
//
// Go Trap:
// Slices, maps, and channels are descriptors/handles that refer to runtime-managed state.
// Copying one does not deep-copy its elements, so aliases can observe mutations.
// Passing `*map[string]int` or `*[]string` creates a redundant pointer-to-pointer,
// forcing clumsy syntax like `(*m)["key"] = val` and confusing callers.
//
// Idiomatic Go:
// Pass slices, maps, and channels by value. Modifications to map entries or slice
// elements in-place are visible through aliases. Appending may replace a slice's
// backing array, so return the resulting slice when the caller needs its new length.
// (Only use `*[]T` if the function needs to reallocate and reassign the caller's slice header).

// UpdateMapAntipattern uses a pointer to a map (unnecessary double indirection).
func UpdateMapAntipattern(m *map[string]int, key string, val int) {
	// Clumsy dereference syntax required
	(*m)[key] = val
}

// UpdateMapIdiomatic passes the map header directly by value.
func UpdateMapIdiomatic(m map[string]int, key string, val int) {
	// Clean, idiomatic access; changes are visible to caller because map is a reference
	m[key] = val
}

// ===========================================================================
// Antipattern 6: Flow Control via Panic and Recover (Using Panic as Exception)
// ===========================================================================
//
// For a Java developer:
// In Java, business validation and ordinary error conditions are thrown as exceptions:
// `throw new InvalidAgeException("must be >= 18")` and caught with `try { ... } catch`.
//
// Go Trap:
// Using `panic()` for business flow control. In Go, `panic()` aborts the goroutine,
// prints a stack trace, and crashes the entire application if unhandled. It is
// computationally expensive and violates Go's explicit error model.
//
// Idiomatic Go:
// "Errors are values."
// Return `(int, error)` and let the caller inspect `if err != nil`.
// Reserve `panic()` strictly for unrecoverable programmer errors (e.g. initialization
// invariant failures, index out of bounds, nil pointer dereferences).

// ParseAgeAntipattern panics on invalid input (Java exception style).
func ParseAgeAntipattern(input string) int {
	age, err := strconv.Atoi(input)
	if err != nil || age < 0 || age > 150 {
		panic(fmt.Sprintf("invalid age input: %s", input))
	}
	return age
}

// ParseAgeIdiomatic returns an error value for invalid input.
func ParseAgeIdiomatic(input string) (int, error) {
	age, err := strconv.Atoi(input)
	if err != nil || age < 0 || age > 150 {
		return 0, fmt.Errorf("invalid age %q: must be between 0 and 150", input)
	}
	return age, nil
}

// ===========================================================================
// Antipattern 7: Slice Memory Retention (Subslice Leak)
// ===========================================================================
//
// For a Java developer:
// This is identical to the classic `String.substring()` memory leak in Java 6, where
// a tiny substring kept the entire large `char[]` backing array alive in the JVM heap.
//
// Go Trap:
// Slicing an existing slice `small := largeData[0:4]` creates a new slice header
// that points to the SAME underlying backing array. As long as `small` is in scope,
// the garbage collector CANNOT reclaim the large backing array (even if it's 100MB+).
//
// Idiomatic Go:
// Allocate a small slice and explicitly copy only the required elements:
// `small := make([]byte, 4); copy(small, largeData[0:4])`. The large slice can now be GC'd.

// ExtractHeaderAntipattern keeps the full 1MB array pinned in memory.
func ExtractHeaderAntipattern(largeData []byte) []byte {
	// BUG: Points to the original largeData backing array; 1MB cannot be garbage collected!
	return largeData[:4]
}

// ExtractHeaderIdiomatic copies the needed bytes to a fresh slice, freeing the original array.
func ExtractHeaderIdiomatic(largeData []byte) []byte {
	header := make([]byte, 4)
	copy(header, largeData[:4])
	return header
}

// ===========================================================================
// Antipattern 8: Slice Growth Without Preallocation
// ===========================================================================
//
// For a Java developer:
// In Java, `new ArrayList<>()` starts with capacity 10 and resizes.
// In Go, growing a slice via `append` without initial capacity causes multiple heap
// allocations and copying data as the backing array grows. The growth strategy is
// an implementation detail; do not depend on a particular capacity sequence.
//
// Idiomatic Go:
// When the target size is known or can be estimated, preallocate capacity:
// `make([]int, 0, capacity)`. This does a single heap allocation.

// GenerateNumbersAntipattern appends to a nil slice, triggering repeated reallocation.
func GenerateNumbersAntipattern(n int) []int {
	var nums []int // capacity 0
	for i := 0; i < n; i++ {
		nums = append(nums, i)
	}
	return nums
}

// GenerateNumbersIdiomatic preallocates capacity to avoid reallocation churn.
func GenerateNumbersIdiomatic(n int) []int {
	nums := make([]int, 0, n) // capacity = n, length = 0
	for i := 0; i < n; i++ {
		nums = append(nums, i)
	}
	return nums
}

// ===========================================================================
// Antipattern 9: Concurrent Map Access (Fatal Runtime Crash)
// ===========================================================================
//
// For a Java developer:
// In Java, accessing a non-thread-safe `HashMap` concurrently might corrupt bucket
// pointers or return null unexpectedly, but it does NOT terminate the JVM process.
//
// Go Trap:
// In Go, concurrent map read and write is detected at the runtime level. When detected,
// Go immediately aborts the process with: `fatal error: concurrent map writes`.
// This fatal panic CANNOT be caught by `recover()`!
//
// Idiomatic Go:
// Protect maps with `sync.RWMutex` (for read-heavy workloads) or `sync.Mutex`,
// or use `sync.Map` for append-only / key-disjoint concurrent caches.

// ThreadSafeMap provides synchronized read/write access to an underlying map.
type ThreadSafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewThreadSafeMap creates a properly initialized ThreadSafeMap.
func NewThreadSafeMap() *ThreadSafeMap {
	return &ThreadSafeMap{
		data: make(map[string]string),
	}
}

// Set acquires an exclusive lock to write a key-value pair.
func (m *ThreadSafeMap) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get acquires a shared read lock to retrieve a value.
func (m *ThreadSafeMap) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// ===========================================================================
// Antipattern 10: Variable Shadowing of 'err'
// ===========================================================================
//
// For a Java developer:
// In Java, you cannot declare two local variables named `err` in the same scope.
//
// Go Trap:
// The short variable declaration `:=` inside an inner block (such as `if` or `for`)
// declares a NEW local variable with the same name, shadowing the outer variable.
// If the outer variable was meant to hold the function's return error, the outer
// error remains `nil`, causing silent failures!
//
// Idiomatic Go:
// Use `=` (assignment) rather than `:=` (short declaration) when targeting an existing
// variable in the enclosing scope, or give inner variables explicit names.

// ParseCoordinatesAntipattern demonstrates how a shadowed error can be discarded.
func ParseCoordinatesAntipattern(sX, sY string) (x, y int, err error) {
	if sX != "" {
		// BUG: ":=" creates a new block-local err. The invalid-input path
		// never assigns the named return err, so the function reports success.
		val, err := strconv.Atoi(sX)
		if err == nil {
			x = val
		}
	}

	if sY != "" {
		val, err := strconv.Atoi(sY)
		if err == nil {
			y = val
		}
	}
	return x, y, nil
}

// ParseCoordinatesIdiomatic propagates conversion errors instead of shadowing them.
func ParseCoordinatesIdiomatic(sX, sY string) (x, y int, err error) {
	if sX != "" {
		x, err = strconv.Atoi(sX)
		if err != nil {
			return 0, 0, fmt.Errorf("parse x coordinate: %w", err)
		}
	}
	if sY != "" {
		y, err = strconv.Atoi(sY)
		if err != nil {
			return 0, 0, fmt.Errorf("parse y coordinate: %w", err)
		}
	}
	return x, y, nil
}

// ProfileSnapshot contains mutable reference-like fields.
type ProfileSnapshot struct {
	Tags       []string
	Attributes map[string]string
}

// ShallowProfileCopyAntipattern copies only the slice and map descriptors.
func ShallowProfileCopyAntipattern(profile ProfileSnapshot) ProfileSnapshot {
	return profile
}

// CloneProfileSnapshot copies the slice and map so mutations do not affect the source.
func CloneProfileSnapshot(profile ProfileSnapshot) ProfileSnapshot {
	profile.Tags = slices.Clone(profile.Tags)
	profile.Attributes = maps.Clone(profile.Attributes)
	return profile
}

// ===========================================================================
// RunAntipatternsDemo
// ===========================================================================

// RunAntipatternsDemo showcases the common Go antipatterns and their idiomatic fixes.
func RunAntipatternsDemo() {
	fmt.Println("--- Go Antipatterns (for Java Developers) Demo ---")

	// 1. The Nil Interface Trap
	fmt.Println("1. The Nil Interface Trap (Typed nil pointer returned as interface):")
	errBad := ValidateUsernameAntipattern("valid_username")
	fmt.Printf("   Antipattern: err != nil is %v (Concrete type: %T, value: %v)\n",
		errBad != nil, errBad, errBad)
	fmt.Println("   -> Calling code sees an error even when validation succeeded!")

	errGood := ValidateUsernameIdiomatic("valid_username")
	fmt.Printf("   Idiomatic:   err != nil is %v (err: %v)\n", errGood != nil, errGood)

	// 2. Passing Mutex by Value
	fmt.Println("2. Passing Mutex by Value (Copies lock state):")
	cVal := ValueCounterAntipattern{}
	cVal = cVal.Inc(1)
	fmt.Printf("   Antipattern: Value receiver copies struct and lock (heldBy inside copy was lost).\n")

	cPtr := &PointerCounterIdiomatic{}
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cPtr.Inc()
		}()
	}
	wg.Wait()
	fmt.Printf("   Idiomatic:   Pointer receiver synchronized count = %d\n", cPtr.Value())

	// 3. Goroutine Leaks
	fmt.Println("3. Goroutine Leak vs Cancellation:")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, errTimeout := QueryDataIdiomatic(ctx, true)
	fmt.Printf("   Idiomatic: Query timed out cleanly: %v (goroutine exited without leak)\n", errTimeout)

	// 4. Consumer-side Interface
	fmt.Println("4. Consumer-Side Interface (Duck Typing vs Producer Interfaces):")
	store := NewCustomerStore()
	name, err := LookupCustomerName(store, "c-1")
	if err != nil {
		fmt.Printf("   Customer lookup failed: %v\n", err)
	} else {
		fmt.Printf("   Lookup Customer: %s (Consumer defined interface 'CustomerFinder')\n", name)
	}

	// 5. Reference Types
	fmt.Println("5. Slices and Maps are Reference Types:")
	m := map[string]int{"alpha": 1}
	UpdateMapIdiomatic(m, "beta", 2)
	fmt.Printf("   Updated map directly by value: %v\n", m)

	// 6. Errors as Values vs Panic
	fmt.Println("6. Errors as Values vs Panic:")
	if age, err := ParseAgeIdiomatic("-5"); err != nil {
		fmt.Printf("   Handled expected error gracefully: %v\n", err)
	} else {
		fmt.Printf("   Parsed age: %d\n", age)
	}

	// 7. Subslice Memory Retention
	fmt.Println("7. Slice Memory Retention (Subslice vs Copy):")
	largeData := make([]byte, 1024*1024) // 1MB
	copy(largeData, []byte("GOGO"))
	retained := ExtractHeaderAntipattern(largeData)
	copied := ExtractHeaderIdiomatic(largeData)
	fmt.Printf("   Retained cap: %d bytes (holds entire array!)\n", cap(retained))
	fmt.Printf("   Copied cap:   %d bytes (allows 1MB array to be garbage collected)\n", cap(copied))

	// 8. Slice Preallocation
	fmt.Println("8. Slice Growth Preallocation:")
	t0 := time.Now()
	_ = GenerateNumbersAntipattern(100000)
	durBad := time.Since(t0)

	t1 := time.Now()
	_ = GenerateNumbersIdiomatic(100000)
	durGood := time.Since(t1)
	fmt.Printf("   Without preallocation: %v | With preallocation: %v\n", durBad, durGood)

	// 9. ThreadSafeMap
	fmt.Println("9. Thread-Safe Map (sync.RWMutex vs Fatal Runtime Crash):")
	safeMap := NewThreadSafeMap()
	safeMap.Set("user_1", "Alice")
	val, _ := safeMap.Get("user_1")
	fmt.Printf("   Retrieved synchronized value: %s\n", val)

	fmt.Println("10. Error Shadowing:")
	_, _, badErr := ParseCoordinatesAntipattern("not-a-number", "20")
	_, _, goodErr := ParseCoordinatesIdiomatic("not-a-number", "20")
	fmt.Printf("   Shadowed error reports failure? %t; explicit propagation reports failure? %t\n",
		badErr != nil, goodErr != nil)

	fmt.Println("11. Shallow Copies of Mutable Fields:")
	profile := ProfileSnapshot{
		Tags:       []string{"go"},
		Attributes: map[string]string{"level": "senior"},
	}
	shallow := ShallowProfileCopyAntipattern(profile)
	shallow.Tags[0] = "java"
	shallow.Attributes["level"] = "staff"
	fmt.Printf("   Shallow copy changed source too: tags=%v attributes=%v\n",
		profile.Tags, profile.Attributes)
	snapshot := CloneProfileSnapshot(profile)
	snapshot.Tags[0] = "isolated"
	snapshot.Attributes["level"] = "principal"
	fmt.Printf("   Deep-enough clone is independent: source=%+v clone=%+v\n", profile, snapshot)

	fmt.Println("--- Go Antipatterns Demo End ---")
}
