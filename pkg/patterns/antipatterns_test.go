package patterns

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// 1. Nil Interface Trap Tests
// ---------------------------------------------------------------------------

func TestNilInterfaceTrap(t *testing.T) {
	// The antipattern: valid input still returns an error interface != nil!
	errBad := ValidateUsernameAntipattern("valid_name")
	if errBad == nil {
		t.Error("Antipattern unexpectedly returned nil; expected non-nil interface holding typed nil pointer")
	}

	// Idiomatic Go: valid input returns true literal nil
	errGood := ValidateUsernameIdiomatic("valid_name")
	if errGood != nil {
		t.Errorf("ValidateUsernameIdiomatic('valid_name') should return nil, got: %v", errGood)
	}

	// Invalid input returns an error in both
	errBadInvalid := ValidateUsernameAntipattern("ab")
	if errBadInvalid == nil {
		t.Error("ValidateUsernameAntipattern('ab') should return error for short name")
	}

	errGoodInvalid := ValidateUsernameIdiomatic("ab")
	if errGoodInvalid == nil {
		t.Error("ValidateUsernameIdiomatic('ab') should return error for short name")
	}
}

// ---------------------------------------------------------------------------
// 2. Mutex Pointer vs Value Tests
// ---------------------------------------------------------------------------

func TestMutexPointerReceiver(t *testing.T) {
	counter := &PointerCounterIdiomatic{}
	const goroutines = 50
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Inc()
			}
		}()
	}
	wg.Wait()

	expected := goroutines * iterations
	if counter.Value() != expected {
		t.Errorf("Expected count %d, got %d", expected, counter.Value())
	}
}

// ---------------------------------------------------------------------------
// 3. Goroutine Leak vs Cancellation Tests
// ---------------------------------------------------------------------------

func TestGoroutineCancellation(t *testing.T) {
	// Fast query succeeds
	ctx := context.Background()
	res, err := QueryDataIdiomatic(ctx, false)
	if err != nil || res != "result-from-slow-query" {
		t.Errorf("Expected success, got res=%q, err=%v", res, err)
	}

	// Slow query with short timeout cancels cleanly without blocking
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, errTimeout := QueryDataIdiomatic(ctxTimeout, true)
	if !errors.Is(errTimeout, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got: %v", errTimeout)
	}
}

// ---------------------------------------------------------------------------
// 4. Consumer-Side Interface & Duck Typing Tests
// ---------------------------------------------------------------------------

type mockCustomerFinder struct {
	mockName string
}

func (m *mockCustomerFinder) GetCustomer(id string) (*Customer, error) {
	if id == "mock-1" {
		return &Customer{ID: id, Name: m.mockName}, nil
	}
	return nil, errors.New("not found")
}

func TestConsumerSideInterface(t *testing.T) {
	// Concrete repository
	store := NewCustomerStore()
	name, err := LookupCustomerName(store, "c-1")
	if err != nil || name != "Acme Corp" {
		t.Errorf("Expected 'Acme Corp', got %q", name)
	}

	// Mock repository satisfies the consumer's CustomerFinder interface seamlessly
	mock := &mockCustomerFinder{mockName: "Mock Enterprise"}
	mockName, err := LookupCustomerName(mock, "mock-1")
	if err != nil || mockName != "Mock Enterprise" {
		t.Errorf("Expected 'Mock Enterprise', got %q", mockName)
	}

	_, err = LookupCustomerName(mock, "unknown-id")
	if err == nil {
		t.Error("Expected lookup error for unknown customer")
	}
}

// ---------------------------------------------------------------------------
// 5. Reference Types Tests
// ---------------------------------------------------------------------------

func TestMapReferenceType(t *testing.T) {
	m := map[string]int{"initial": 10}
	UpdateMapIdiomatic(m, "added", 20)

	if m["added"] != 20 {
		t.Errorf("Expected m['added'] = 20, got %d", m["added"])
	}

	// Pointer-to-map antipattern achieves the same but with clunky syntax
	UpdateMapAntipattern(&m, "antipattern", 30)
	if m["antipattern"] != 30 {
		t.Errorf("Expected m['antipattern'] = 30, got %d", m["antipattern"])
	}
}

// ---------------------------------------------------------------------------
// 6. Errors as Values vs Panic Tests
// ---------------------------------------------------------------------------

func TestParseAge(t *testing.T) {
	// Idiomatic success
	age, err := ParseAgeIdiomatic("25")
	if err != nil || age != 25 {
		t.Errorf("ParseAgeIdiomatic('25') = (%d, %v), want (25, nil)", age, err)
	}

	// Idiomatic failure returns error value, does not panic
	_, err = ParseAgeIdiomatic("-1")
	if err == nil {
		t.Error("ParseAgeIdiomatic('-1') expected error, got nil")
	}

	// Antipattern panics on invalid input
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("ParseAgeAntipattern('-1') should have panicked")
			}
		}()
		ParseAgeAntipattern("-1")
	}()
}

// ---------------------------------------------------------------------------
// 7. Subslice Memory Retention Tests
// ---------------------------------------------------------------------------

func TestSubsliceMemoryRetention(t *testing.T) {
	const bigSize = 10000
	largeData := make([]byte, bigSize)

	retained := ExtractHeaderAntipattern(largeData)
	if cap(retained) != bigSize {
		t.Errorf("Antipattern slice cap = %d, expected entire backing array (%d)", cap(retained), bigSize)
	}

	copied := ExtractHeaderIdiomatic(largeData)
	if cap(copied) != 4 {
		t.Errorf("Idiomatic slice cap = %d, expected exactly 4", cap(copied))
	}
}

// ---------------------------------------------------------------------------
// 8. Slice Preallocation Tests
// ---------------------------------------------------------------------------

func TestSlicePreallocation(t *testing.T) {
	bad := GenerateNumbersAntipattern(50)
	good := GenerateNumbersIdiomatic(50)

	if !reflect.DeepEqual(bad, good) {
		t.Errorf("Expected slices to match: bad=%v, good=%v", bad, good)
	}
	if len(good) != 50 || cap(good) != 50 {
		t.Errorf("Idiomatic slice length=%d, cap=%d, expected 50 and 50", len(good), cap(good))
	}
}

// ---------------------------------------------------------------------------
// 9. ThreadSafeMap Tests
// ---------------------------------------------------------------------------

func TestThreadSafeMap(t *testing.T) {
	tsMap := NewThreadSafeMap()
	var wg sync.WaitGroup

	// Concurrently write and read keys
	for i := 0; i < 20; i++ {
		wg.Add(2)
		key := "key"
		val := "val"
		go func() {
			defer wg.Done()
			tsMap.Set(key, val)
		}()
		go func() {
			defer wg.Done()
			_, _ = tsMap.Get(key)
		}()
	}
	wg.Wait()

	v, ok := tsMap.Get("key")
	if !ok || v != "val" {
		t.Errorf("Expected key 'key'='val', got ok=%v, v=%s", ok, v)
	}
}

// ---------------------------------------------------------------------------
// 10. Variable Shadowing Tests
// ---------------------------------------------------------------------------

func TestVariableShadowing(t *testing.T) {
	x, y, err := ParseCoordinatesAntipattern("10", "20")
	if err != nil || x != 10 || y != 20 {
		t.Errorf("Expected (10, 20, nil), got (%d, %d, %v)", x, y, err)
	}

	// The block-local err is discarded, so the antipattern incorrectly reports success.
	x, y, err = ParseCoordinatesAntipattern("invalid", "20")
	if err != nil || x != 0 || y != 20 {
		t.Errorf("Antipattern should hide the invalid X error, got (%d, %d, %v)", x, y, err)
	}

	_, _, err = ParseCoordinatesIdiomatic("invalid", "20")
	if err == nil {
		t.Error("Idiomatic parser should return an error for an invalid X coordinate")
	}

	x, y, err = ParseCoordinatesIdiomatic("10", "20")
	if err != nil || x != 10 || y != 20 {
		t.Errorf("Idiomatic parser = (%d, %d, %v), want (10, 20, nil)", x, y, err)
	}
}

func TestShallowCopyAndCloneOfMutableFields(t *testing.T) {
	profile := ProfileSnapshot{
		Tags:       []string{"go"},
		Attributes: map[string]string{"level": "senior"},
	}

	shallow := ShallowProfileCopyAntipattern(profile)
	shallow.Tags[0] = "java"
	shallow.Attributes["level"] = "staff"
	if profile.Tags[0] != "java" || profile.Attributes["level"] != "staff" {
		t.Fatal("shallow copy should share the slice backing array and map entries")
	}

	clone := CloneProfileSnapshot(profile)
	clone.Tags[0] = "rust"
	clone.Attributes["level"] = "principal"
	if profile.Tags[0] != "java" || profile.Attributes["level"] != "staff" {
		t.Fatal("clone mutations should not affect the original profile")
	}

	if got := CloneProfileSnapshot(ProfileSnapshot{}); got.Tags != nil || got.Attributes != nil {
		t.Fatal("cloning nil fields should preserve nil")
	}
}
