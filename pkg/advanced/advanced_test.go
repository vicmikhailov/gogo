package advanced // Declaring the package as advanced for testing.

import ( // Starting the import block for test dependencies.
	"context" // Importing context for testing context-aware functions.
	"math"    // Importing math for mathematical operations like Abs.
	"reflect" // Importing reflect for DeepEqual comparisons.
	"strings" // Importing strings for string content checks.
	"sync"    // Importing sync for WaitGroup in concurrency tests.
	"testing" // Importing the testing package for Go's test framework.
	"time"    // Importing time for duration and timeout testing.
) // Closing the import block.

/**
 * ===========================================================================
 * Advanced Feature Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Go's testing style is often more "table-driven" (see TestSafeDivide for example).
 * - There is no `@Before` or `@After` annotation by default; use `t.Cleanup()`
 *   or just defer at the beginning of the test.
 * - Test state is usually managed locally within the test function.
 */ // Block comment explaining Go testing to Java developers.

// ---------------------------------------------------------------------------
// Embedding // Section header for embedding tests.
// ---------------------------------------------------------------------------

func TestDogEmbedding(t *testing.T) { // Test function for Dog embedding.
	dog := Dog{Animal: Animal{Name: "Rex", Legs: 4}, Breed: "Labrador"} // Initializing a Dog with embedded Animal.
	// Promoted field access // Testing access to fields from the embedded struct.
	if dog.Name != "Rex" { // Checking if Name field was promoted.
		t.Errorf("Expected Name 'Rex', got %s", dog.Name) // Logging error if Name is incorrect.
	} // End of Name check.
	if dog.Legs != 4 { // Checking if Legs field was promoted.
		t.Errorf("Expected Legs 4, got %d", dog.Legs) // Logging error if Legs is incorrect.
	} // End of Legs check.
	// Promoted method // Testing access to methods from the embedded struct.
	desc := dog.Describe()              // Calling the promoted Describe method.
	if !strings.Contains(desc, "Rex") { // Checking if the description contains the Name.
		t.Errorf("Describe() should contain 'Rex', got %s", desc) // Logging error if description is incorrect.
	} // End of Describe check.
	// Own method // Testing the struct's own method.
	speak := dog.Speak()                  // Calling the Speak method.
	if !strings.Contains(speak, "Woof") { // Checking if the speech contains "Woof".
		t.Errorf("Speak() should contain 'Woof', got %s", speak) // Logging error if speech is incorrect.
	} // End of Speak check.
} // Closing the TestDogEmbedding function.

func TestServiceDogOverride(t *testing.T) { // Test function for ServiceDog method overriding.
	sd := ServiceDog{ // Initializing a ServiceDog.
		Dog:    Dog{Animal: Animal{Name: "Buddy", Legs: 4}, Breed: "Golden"}, // Embedding a Dog.
		CertID: "SD-99",                                                      // Setting the CertID.
	} // End of initialization.
	desc := sd.Describe()                                                           // Calling the overridden Describe method.
	if !strings.Contains(desc, "Service Dog") || !strings.Contains(desc, "SD-99") { // Checking for expected content.
		t.Errorf("ServiceDog.Describe() should mention cert, got %s", desc) // Logging error if description is incorrect.
	} // End of Describe check.
} // Closing the TestServiceDogOverride function.

// ---------------------------------------------------------------------------
// Enum (iota) // Section header for enum tests.
// ---------------------------------------------------------------------------

func TestColorEnum(t *testing.T) { // Test function for Color enum.
	if Red.String() != "Red" { // Checking the String representation of Red.
		t.Errorf("Expected 'Red', got %s", Red.String()) // Logging error if string is incorrect.
	} // End of Red.String check.
	if !Red.IsWarm() { // Checking if Red is correctly identified as warm.
		t.Error("Red should be warm") // Logging error if Red is not warm.
	} // End of Red.IsWarm check.
	if Blue.IsWarm() { // Checking if Blue is correctly identified as NOT warm.
		t.Error("Blue should not be warm") // Logging error if Blue is warm.
	} // End of Blue.IsWarm check.
	all := AllColors() // Getting all colors from the utility function.
	if len(all) != 4 { // Checking the number of colors.
		t.Errorf("Expected 4 colors, got %d", len(all)) // Logging error if count is incorrect.
	} // End of count check.
} // Closing the TestColorEnum function.

// ---------------------------------------------------------------------------
// Functional patterns // Section header for functional pattern tests.
// ---------------------------------------------------------------------------

func TestCompose(t *testing.T) { // Test function for function composition.
	double := func(x int) int { return x * 2 } // Defining a double function.
	addOne := func(x int) int { return x + 1 } // Defining an increment function.
	// Compose(addOne, double) means addOne(double(x)) // Explaining the composition order.
	fn := Compose(addOne, double) // Composing the two functions.
	if fn(3) != 7 {               // Checking the result of composition (3*2+1).
		t.Errorf("Expected 7, got %d", fn(3)) // Logging error if result is incorrect.
	} // End of result check.
} // Closing the TestCompose function.

func TestCurry2(t *testing.T) { // Test function for function currying.
	add := func(a, b int) int { return a + b } // Defining a binary addition function.
	add5 := Curry2(add)(5)                     // Currying the addition function and applying the first argument.
	if add5(3) != 8 {                          // Checking the result of the curried function (5+3).
		t.Errorf("Expected 8, got %d", add5(3)) // Logging error if result is incorrect.
	} // End of result check.
} // Closing the TestCurry2 function.

func TestMemoize(t *testing.T) { // Test function for function memoization.
	calls := 0                      // Counter for function computations.
	fn := Memoize(func(n int) int { // Creating a memoized squaring function.
		calls++      // Incrementing the counter on each actual computation.
		return n * n // Returning the square.
	}) // End of memoized function.
	fn(3)           // First call with key 3.
	fn(3)           // Second call with key 3.
	fn(3)           // Third call with key 3.
	if calls != 1 { // Checking that only one computation was performed.
		t.Errorf("Expected 1 computation, got %d", calls) // Logging error if more than one computation occurred.
	} // End of calls check.
	if fn(4) != 16 { // Calling with a new key 4.
		t.Errorf("Expected 16, got %d", fn(4)) // Logging error if result is incorrect.
	} // End of result check.
	if calls != 2 { // Checking that a second computation was performed for the new key.
		t.Errorf("Expected 2 computations after new key, got %d", calls) // Logging error if count is incorrect.
	} // End of total calls check.
} // Closing the TestMemoize function.

func TestPipeline(t *testing.T) { // Test function for the functional pipeline.
	result := Pipeline("  hello  ", // Executing the pipeline on a string.
		strings.TrimSpace, // Stage 1: trimming whitespace.
		strings.ToUpper,   // Stage 2: converting to uppercase.
	) // End of pipeline.
	if result != "HELLO" { // Checking the final result.
		t.Errorf("Expected 'HELLO', got %s", result) // Logging error if result is incorrect.
	} // End of result check.
} // Closing the TestPipeline function.

// ---------------------------------------------------------------------------
// Defer / Panic / Recover // Section header for Defer/Panic/Recover tests.
// ---------------------------------------------------------------------------

func TestSafeDivide(t *testing.T) { // Test function for SafeDivide using table-driven tests.
	tests := []struct { // Defining the test table structure.
		name      string // Name of the test case.
		a, b      int    // Input values.
		want      int    // Expected result value.
		wantError bool   // Whether an error is expected.
	}{ // Starting the test cases.
		{"normal division", 10, 2, 5, false},   // Case: valid division.
		{"divide by zero", 10, 0, 0, true},     // Case: division by zero.
		{"negative result", -10, 2, -5, false}, // Case: negative result.
	} // End of test cases.

	for _, tt := range tests { // Iterating through each test case.
		t.Run(tt.name, func(t *testing.T) { // Running a subtest for each case.
			got, err := SafeDivide(tt.a, tt.b) // Calling the function being tested.
			if (err != nil) != tt.wantError {  // Checking if the error status matches expectation.
				t.Errorf("SafeDivide() error = %v, wantError %v", err, tt.wantError) // Logging error if mismatch.
				return                                                               // Exiting subtest on error status mismatch.
			} // End of error status check.
			if got != tt.want { // Checking if the result value matches expectation.
				t.Errorf("SafeDivide() got = %v, want %v", got, tt.want) // Logging error if result is incorrect.
			} // End of value check.
		}) // End of subtest.
	} // End of test cases loop.
} // Closing the TestSafeDivide function.

// ---------------------------------------------------------------------------
// Reflection // Section header for reflection tests.
// ---------------------------------------------------------------------------

func TestDescribeType(t *testing.T) { // Test function for reflection-based type description.
	type Sample struct { // Defining a sample struct for testing.
		X int    // Field X of type int.
		Y string // Field Y of type string.
	} // Closing the Sample struct definition.
	desc := DescribeType(Sample{X: 42, Y: "hello"}) // Getting the description for a Sample instance.
	if !strings.Contains(desc, "Sample") {          // Checking if the type name is in the description.
		t.Errorf("Expected type name 'Sample' in description, got %s", desc) // Logging error if missing.
	} // End of name check.
	if !strings.Contains(desc, "Fields: 2") { // Checking if the number of fields is correct.
		t.Errorf("Expected 'Fields: 2', got %s", desc) // Logging error if count is incorrect.
	} // End of fields check.
} // Closing the TestDescribeType function.

func TestStructToMap(t *testing.T) { // Test function for reflection-based struct-to-map conversion.
	type Sample struct { // Defining a sample struct for testing.
		Name string // Field Name.
		Age  int    // Field Age.
	} // Closing the Sample struct definition.
	m := StructToMap(Sample{Name: "Alice", Age: 30}) // Converting Sample instance to a map.
	if m["Name"] != "Alice" {                        // Checking if Name is correctly mapped.
		t.Errorf("Expected Name='Alice', got %v", m["Name"]) // Logging error if incorrect.
	} // End of Name check.
	if m["Age"] != 30 { // Checking if Age is correctly mapped.
		t.Errorf("Expected Age=30, got %v", m["Age"]) // Logging error if incorrect.
	} // End of Age check.
} // Closing the TestStructToMap function.

// ---------------------------------------------------------------------------
// Type constraints (Number) // Section header for generic constraint tests.
// ---------------------------------------------------------------------------

func TestSum(t *testing.T) { // Test function for the generic Sum function.
	if Sum([]int{1, 2, 3, 4}) != 10 { // Testing Sum with a slice of integers.
		t.Errorf("Expected 10, got %d", Sum([]int{1, 2, 3, 4})) // Logging error if sum is incorrect.
	} // End of int sum check.
	if math.Abs(Sum([]float64{1.5, 2.5})-4.0) > 1e-9 { // Testing Sum with a slice of floats.
		t.Errorf("Expected 4.0, got %f", Sum([]float64{1.5, 2.5})) // Logging error if sum is incorrect.
	} // End of float sum check.
} // Closing the TestSum function.

func TestMinMax(t *testing.T) { // Test function for generic Min and Max functions.
	items := []int{3, 1, 4, 1, 5, 9} // Initializing a slice of integers.
	if Min(items) != 1 {             // Checking the minimum value.
		t.Errorf("Expected Min=1, got %d", Min(items)) // Logging error if Min is incorrect.
	} // End of Min check.
	if Max(items) != 9 { // Checking the maximum value.
		t.Errorf("Expected Max=9, got %d", Max(items)) // Logging error if Max is incorrect.
	} // End of Max check.
} // Closing the TestMinMax function.

func TestClamp(t *testing.T) { // Test function for the generic Clamp function.
	tests := []struct { // Defining the test table structure.
		val, lo, hi int // Input value and bounds.
		want        int // Expected clamped result.
	}{ // Starting test cases.
		{15, 0, 10, 10}, // Case: value above upper bound.
		{-5, 0, 10, 0},  // Case: value below lower bound.
		{5, 0, 10, 5},   // Case: value within bounds.
	} // End of test cases.
	for _, tt := range tests { // Iterating through test cases.
		got := Clamp(tt.val, tt.lo, tt.hi) // Calling the function being tested.
		if got != tt.want {                // Checking if the result matches expectation.
			t.Errorf("Clamp(%d, %d, %d) = %d; want %d", tt.val, tt.lo, tt.hi, got, tt.want) // Logging error if incorrect.
		} // End of result check.
	} // End of loop.
} // Closing the TestClamp function.

// ---------------------------------------------------------------------------
// FanOut // Section header for the FanOut concurrency test.
// ---------------------------------------------------------------------------

func TestFanOut(t *testing.T) { // Test function for the FanOut generic concurrent processor.
	inputs := []int{1, 2, 3, 4, 5}                                 // Initializing input slice.
	results := FanOut(inputs, 2, func(n int) int { return n * n }) // Executing FanOut with 2 workers to square numbers.
	expected := []int{1, 4, 9, 16, 25}                             // Defining the expected squared results.
	if !reflect.DeepEqual(results, expected) {                     // Using DeepEqual to compare slices.
		t.Errorf("Expected %v, got %v", expected, results) // Logging error if results are incorrect.
	} // End of results check.
} // Closing the TestFanOut function.

// ---------------------------------------------------------------------------
// Java FAQ Tests // Section header for Java FAQ related tests.
// ---------------------------------------------------------------------------

func TestAccount(t *testing.T) { // Test function for the Account struct and encapsulation.
	acc := Account{Owner: "Alice"} // Initializing an Account.
	acc.Deposit(100)               // Depositing 100.
	if acc.GetBalance() != 100 {   // Checking the balance via getter.
		t.Errorf("Expected balance 100, got %d", acc.GetBalance()) // Logging error if balance is incorrect.
	} // End of balance check.
	acc.Deposit(-50)             // Attempting a negative deposit.
	if acc.GetBalance() != 100 { // Checking that the negative deposit was ignored.
		t.Errorf("Expected balance 100 after negative deposit, got %d", acc.GetBalance()) // Logging error if balance changed.
	} // End of negative deposit check.
} // Closing the TestAccount function.

func TestUserRecord(t *testing.T) { // Test function for UserRecord data transfer object.
	u := UserRecord{Username: "gopher", Email: "go@golang.org"} // Initializing a UserRecord.
	if u.Username != "gopher" {                                 // Checking the Username field.
		t.Errorf("Expected Username 'gopher', got %s", u.Username) // Logging error if Username is incorrect.
	} // End of Username check.
} // Closing the TestUserRecord function.

// ---------------------------------------------------------------------------
// Advanced Samples Tests // Section header for various advanced sample tests.
// ---------------------------------------------------------------------------

func TestDatabaseOptions(t *testing.T) { // Test function for the functional options pattern.
	db := NewDatabaseConnector( // Creating a connector with specific options.
		WithHost("example.com"),    // Setting the host.
		WithPort(9999),             // Setting the port.
		WithTimeout(5*time.Second), // Setting the timeout.
	) // End of connector creation.
	if db.Host != "example.com" { // Checking the Host field.
		t.Errorf("Expected Host 'example.com', got %s", db.Host) // Logging error if Host is incorrect.
	} // End of Host check.
	if db.Port != 9999 { // Checking the Port field.
		t.Errorf("Expected Port 9999, got %d", db.Port) // Logging error if Port is incorrect.
	} // End of Port check.
	if db.Timeout != 5*time.Second { // Checking the Timeout field.
		t.Errorf("Expected Timeout 5s, got %v", db.Timeout) // Logging error if Timeout is incorrect.
	} // End of Timeout check.
} // Closing the TestDatabaseOptions function.

func TestSafeCounter(t *testing.T) { // Test function for the atomic SafeCounter.
	counter := &SafeCounter{}   // Initializing a SafeCounter.
	var wg sync.WaitGroup       // Declaring a WaitGroup for concurrency.
	for i := 0; i < 1000; i++ { // Launching 1000 concurrent increments.
		wg.Add(1)   // Incrementing WaitGroup counter.
		go func() { // Launching a goroutine.
			defer wg.Done()     // Ensuring Done is called.
			counter.Increment() // Incrementing the counter.
		}() // End of goroutine.
	} // End of loop.
	wg.Wait()                    // Waiting for all increments to complete.
	if counter.Value() != 1000 { // Checking the final counter value.
		t.Errorf("Expected Value 1000, got %d", counter.Value()) // Logging error if count is incorrect.
	} // End of value check.
} // Closing the TestSafeCounter function.

func TestConcurrentCache(t *testing.T) { // Test function for the generic concurrent cache.
	cache := NewConcurrentCache[string, string]() // Initializing a new cache.
	cache.Set("key1", "value1")                   // Setting a key-value pair.
	val, ok := cache.Get("key1")                  // Retrieving the value.
	if !ok || val != "value1" {                   // Checking the retrieved value and existence status.
		t.Errorf("Expected value1, got %v, ok: %v", val, ok) // Logging error if incorrect.
	} // End of retrieval check.
	_, ok = cache.Get("key2") // Attempting to retrieve a non-existent key.
	if ok {                   // Checking that it was not found.
		t.Error("Expected key2 not to be in cache") // Logging error if found.
	} // End of non-existent key check.
} // Closing the TestConcurrentCache function.

func TestAdvancedPipeline(t *testing.T) { // Test function for the advanced concurrent pipeline.
	ctx := context.Background()                          // Creating a background context.
	sum, err := AdvancedPipeline(ctx, []int{1, 2, 3, 4}) // Executing the pipeline.
	if err != nil {                                      // Checking for pipeline errors.
		t.Fatalf("Expected no error, got %v", err) // Failing test if error occurred.
	} // End of error check.
	// Squares: 1, 4, 9, 16. Evens: 4, 16. Sum: 20. // Explaining the expected calculation.
	if sum != 20 { // Checking the final sum result.
		t.Errorf("Expected sum 20, got %d", sum) // Logging error if sum is incorrect.
	} // End of sum check.

	// Test cancellation // Section for testing pipeline cancellation.
	ctx, cancel := context.WithCancel(context.Background()) // Creating a cancellable context.
	cancel()                                                // Cancelling the context immediately.
	_, err = AdvancedPipeline(ctx, []int{1, 2, 3})          // Attempting to run the pipeline with cancelled context.
	if err == nil {                                         // Checking that an error was returned.
		t.Error("Expected error from cancelled context, got nil") // Logging error if no error returned.
	} // End of cancellation check.
} // Closing the TestAdvancedPipeline function.
