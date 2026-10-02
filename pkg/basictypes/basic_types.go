// Package basictypes showcases Go's built-in types (slices, maps, strings, JSON).
//
// For a Java developer:
//   - Go does not have a universal `null`. It has `nil` which is the zero-value
//     for pointers, interfaces, maps, slices, and channels.
//   - Java-ism to avoid: Checking `if (x == nil)` for a primitive type (int, bool, string).
//     In Go, these types cannot be `nil`. They have their own zero-values (0, false, "").
//   - Strings are immutable byte sequences (UTF-8 by default), not `char[]`.
package basictypes

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// ---------------------------------------------------------------------------
// 1. Slice (List) Manipulation (Standard idiomatic ways)
// ---------------------------------------------------------------------------

// RunSliceManipulationDemo demonstrates standard slice operations.
//
// For a Java developer:
//   - Go's `slice` is a dynamic array view (like `ArrayList<T>`).
//   - Slices are views into an underlying array. When you re-slice, they share memory.
//   - The zero-value of a slice is `nil` (like an uninitialized Java List).
//   - There is no formal `List` interface in the standard library.
//   - A slice is a small descriptor (backing array, length, capacity) passed by value.
//   - Operations:
//   - Creation: make([]T, len, cap)  (Java: new ArrayList<>(cap))
//   - Adding:   append(slice, items) (Java: list.add(item))
//   - Slicing:  slice[start:end]     (Java: list.subList(start, end) - half-open interval)
//   - Capacity: len() is current size, cap() is underlying buffer size
//   - Range:    for i, v := range slice (Java: for loop or Stream with index)
//   - Copying:  copy(dest, src) copies elements between slices
func RunSliceManipulationDemo() {
	fmt.Println("\n--- Slice (List) Manipulation Demo ---")

	// a. Creation: make(type, len, cap)
	// Java equivalent: List<Integer> list = new ArrayList<>(10);
	nums := make([]int, 0, 5)
	fmt.Printf("   Initial: len=%d, cap=%d, %v\n", len(nums), cap(nums), nums)

	// b. Adding: append(slice, element...)
	// Java equivalent: list.add(10);
	nums = append(nums, 1, 2, 3)
	fmt.Printf("   After append: %v\n", nums)

	// c. Slicing: slice[start:end] (half-open interval [start, end))
	// Java equivalent: list.subList(1, 3);
	sub := nums[1:3]
	fmt.Printf("   Sub-slice [1:3]: %v\n", sub)

	// d. Length and Capacity
	fmt.Printf("   Length: %d, Capacity: %d\n", len(nums), cap(nums))

	// e. Iteration: range returns index and value
	// Java equivalent: for (int i=0; i < list.size(); i++) { ... }
	fmt.Print("   Iteration: ")
	for i, v := range nums {
		fmt.Printf("[%d]:%d ", i, v)
	}
	fmt.Println()

	// f. Copying: copy(dest, src)
	backup := make([]int, len(nums))
	copy(backup, nums)
	fmt.Printf("   Copy (backup): %v\n", backup)
}

// ---------------------------------------------------------------------------
// 2. Map Manipulation (Standard idiomatic ways)
// ---------------------------------------------------------------------------

// RunMapManipulationDemo demonstrates hash map operations in Go.
//
// For a Java developer:
//   - Go's `map` is a hash map (like `HashMap<K, V>`).
//   - Iteration order is unspecified; sort keys when output must be deterministic.
//   - Copying a map value shares its underlying entries, so mutations are visible through aliases.
//   - Accessing a non-existent key returns the zero-value (0, "", nil) instead of throwing an exception or returning null.
//   - The "comma ok" idiom (`v, ok := m[key]`) is used to test whether a key exists.
//   - To iterate in deterministic order, extract keys into a slice and sort them first.
func RunMapManipulationDemo() {
	fmt.Println("\n--- Map Manipulation Demo ---")

	// a. Creation: make(map[KeyType]ValueType)
	// Java equivalent: Map<String, Integer> map = new HashMap<>();
	ages := make(map[string]int)

	// b. Adding / Updating
	// Java equivalent: map.put("Alice", 30);
	ages["Alice"] = 30
	ages["Bob"] = 25
	fmt.Printf("   Initial map: %v\n", ages)

	// c. Existence check: the "comma ok" idiom
	// Java equivalent: map.containsKey("Alice");
	age, ok := ages["Alice"]
	if ok {
		fmt.Printf("   Alice's age is %d\n", age)
	}

	// d. Deleting: delete(map, key)
	// Java equivalent: map.remove("Bob");
	delete(ages, "Bob")
	fmt.Printf("   After delete(Bob): %v\n", ages)

	// e. Iteration (Warning: Order is unspecified!)
	fmt.Print("   Iteration (order varies): ")
	for name, age := range ages {
		fmt.Printf("%s:%d ", name, age)
	}
	fmt.Println()

	// f. Deterministic Iteration (Sort keys first)
	fmt.Print("   Deterministic iteration: ")
	ages["Charlie"] = 35
	ages["Alpha"] = 20
	keys := make([]string, 0, len(ages))
	for k := range ages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s:%d ", k, ages[k])
	}
	fmt.Println()

	// g. A map can implement a set, and its zero-value element makes counters concise.
	// Java equivalents: Set<String> and Map<String, Integer>.merge(key, 1, Integer::sum).
	seen := make(map[string]struct{})
	wordCounts := make(map[string]int)
	for _, word := range []string{"go", "maps", "go", "sets"} {
		seen[word] = struct{}{}
		wordCounts[word]++
	}
	_, hasGo := seen["go"]
	fmt.Printf("   Set membership: go=%t; frequency map: %v\n", hasGo, wordCounts)

	// h. Nested maps model grouped data; initialize each inner map before writing.
	teams := map[string]map[string]int{}
	if teams["platform"] == nil {
		teams["platform"] = make(map[string]int)
	}
	teams["platform"]["gophers"] = 4
	fmt.Printf("   Nested map: %v\n", teams)
}

// RunMutabilityDemo contrasts Go's value copies with shared slice/map storage.
//
// For a Java developer:
//   - Struct assignment copies fields; it does not create an alias to the original struct.
//   - A slice is a copied descriptor (array pointer, length, capacity), so element writes
//     are shared when descriptors refer to the same backing array.
//   - Appending returns a possibly new descriptor; callers must use the returned slice.
//   - A map assignment copies a descriptor to shared map storage; entry updates are visible
//     through every copy. None of these types make concurrent mutation safe.
func RunMutabilityDemo() {
	fmt.Println("\n--- Mutability and Copy Semantics Demo ---")

	values := make([]int, 3, 4)
	copy(values, []int{1, 2, 3})
	alias := values
	alias[0] = 99
	fmt.Printf("   Slice element mutation is shared: values=%v alias=%v\n", values, alias)

	grown := append(values, 4)
	fmt.Printf("   Append returns a new slice header: old len=%d, new=%v\n", len(values), grown)

	cloned := append([]int(nil), values...)
	cloned[0] = 7
	fmt.Printf("   Explicit copy is independent: values=%v clone=%v\n", values, cloned)

	original := struct {
		Name string
	}{Name: "before"}
	valueCopy := original
	valueCopy.Name = "value copy"
	pointer := &original
	pointer.Name = "pointer mutation"
	fmt.Printf("   Struct assignment copies fields: original=%q copy=%q\n", original.Name, valueCopy.Name)

	metadata := map[string]string{"owner": "platform"}
	metadataAlias := metadata
	metadataAlias["owner"] = "runtime"
	fmt.Printf("   Map entry mutation is shared: %v\n", metadata)

	fmt.Println("--- Mutability and Copy Semantics End ---")
}

// ---------------------------------------------------------------------------
// 3. String Operations (Standard Library)
// ---------------------------------------------------------------------------

// RunStringOperationsDemo showcases string functions from the standard library.
//
// For a Java developer:
//   - Strings in Go are immutable sequences of bytes (UTF-8 encoded by default).
//   - Unlike Java's `String` (which is backed by `char[]` / UTF-16), Go strings are `byte[]`.
//   - `len(str)` returns the number of BYTES, not characters.
//   - Use `[]rune(str)` to get Unicode code points (runes).
//   - Use `strings.Builder` for efficient concatenation in loops (like Java `StringBuilder`).
//   - Go backticks (`) define raw multiline strings (like Java 15+ Text Blocks).
func RunStringOperationsDemo() {
	fmt.Println("\n--- String Operations Demo ---")

	text := "Go is a statically typed, compiled programming language."

	// a. Basic checks: Contains, HasPrefix, HasSuffix
	fmt.Printf("   Contains 'typed':   %t\n", strings.Contains(text, "typed"))
	fmt.Printf("   Has prefix 'Go':    %t\n", strings.HasPrefix(text, "Go"))
	fmt.Printf("   Has suffix 'Java':  %t\n", strings.HasSuffix(text, "Java"))

	// b. Manipulation: ToUpper, Replace
	fmt.Printf("   Upper:              %s\n", strings.ToUpper("go rocks"))
	fmt.Printf("   Replace:            %s\n", strings.Replace(text, "Go", "Golang", 1))

	// c. Splitting and Joining: Fields, Join
	words := strings.Fields(text)
	fmt.Printf("   Words count:        %d\n", len(words))
	joined := strings.Join(words[:3], "-")
	fmt.Printf("   Join first three:   %s\n", joined)

	// d. Trimming whitespace
	dirty := "   \t hello world \n  "
	fmt.Printf("   Trimmed:           '%s'\n", strings.TrimSpace(dirty))

	// e. strings.Builder (Performance like StringBuilder)
	var builder strings.Builder
	for i := 1; i <= 3; i++ {
		builder.WriteString(fmt.Sprintf("Step %d; ", i))
	}
	fmt.Printf("   Builder result:     %s\n", builder.String())

	// f. Unicode / Runes
	// Java equivalent: String.codePointAt()
	japanese := "こんにちは"
	fmt.Printf("   Bytes length:       %d (not characters!)\n", len(japanese))
	fmt.Printf("   Runes count:        %d\n", len([]rune(japanese)))
}

// ---------------------------------------------------------------------------
// 4. JSON Manipulation (Standard Library)
// ---------------------------------------------------------------------------

// Product demonstrates struct tags used for JSON serialization.
//
// For a Java developer:
//   - Struct tags like `json:"id"` are similar to Jackson's `@JsonProperty("id")`.
//   - Fields must be uppercase (Exported) for `encoding/json` to access them.
//   - `omitempty` omits the field when serializing if it has its zero value.
//   - `json:"-"` ignores the field entirely during JSON serialization.
type Product struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Price float64  `json:"price"`
	Tags  []string `json:"tags,omitempty"`
}

// RunJSONDemo demonstrates JSON marshaling and unmarshaling.
//
// For a Java developer:
//   - Marshalling = Serializing (Object to JSON, like `objectMapper.writeValueAsString()`).
//   - Unmarshalling = Deserializing (JSON to Object, like `objectMapper.readValue()`).
//   - Error handling is explicit via return values, unlike Java's `JsonProcessingException`.
//   - Dynamic schemas can be deserialized into `map[string]any` (like `Map<String, Object>`).
func RunJSONDemo() {
	fmt.Println("\n--- JSON Manipulation Demo ---")

	// a. Marshalling (Struct to JSON)
	p := Product{
		ID:    101,
		Name:  "Gopher Plushie",
		Price: 19.99,
		Tags:  []string{"toy", "mascot"},
	}
	jsonData, err := json.MarshalIndent(p, "   ", "  ")
	if err != nil {
		fmt.Printf("   Could not marshal product: %v\n", err)
		return
	}
	fmt.Printf("   JSON Output:\n%s\n", string(jsonData))

	// b. Unmarshalling (JSON to Struct)
	rawJSON := `{"id": 102, "name": "Go Mug", "price": 12.50}`
	var p2 Product
	if err := json.Unmarshal([]byte(rawJSON), &p2); err != nil {
		fmt.Printf("   Could not unmarshal product: %v\n", err)
		return
	}
	fmt.Printf("   Unmarshalled Struct: %+v\n", p2)

	// c. Arbitrary JSON (using map[string]any)
	// Useful when the schema is dynamic or unknown (like Java Map<String, Object>).
	var data map[string]any
	if err := json.Unmarshal([]byte(rawJSON), &data); err != nil {
		fmt.Printf("   Could not unmarshal dynamic JSON: %v\n", err)
		return
	}
	fmt.Printf("   Map representation:  %v (Name: %v)\n", data, data["name"])

	// d. JSON with Custom Logic (omitempty)
	pEmpty := Product{ID: 1, Name: "Invisible"}
	emptyJSON, err := json.Marshal(pEmpty)
	if err != nil {
		fmt.Printf("   Could not marshal empty product: %v\n", err)
		return
	}
	fmt.Printf("   Omitempty Tags:     %s\n", string(emptyJSON))
}

// RunBasicTypesDemo is the facade orchestrator for all basic types demonstrations.
func RunBasicTypesDemo() {
	fmt.Println("--- Basic Types & Standard Library Demo ---")
	RunSliceManipulationDemo()
	RunMapManipulationDemo()
	RunMutabilityDemo()
	RunStringOperationsDemo()
	RunJSONDemo()
	fmt.Println("--- Basic Types & Standard Library End ---")
}

// ReverseString reverses a UTF-8 string by operating on runes.
//
// For a Java developer:
//   - Strings in Go are UTF-8 encoded by default.
//   - Iterating over a string using `range` yields `runes` (Unicode code points).
//   - `rune` is an alias for `int32`.
//   - Operating on runes ensures multi-byte characters (e.g. Japanese, emojis) are not corrupted.
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome reports whether s is a palindrome, ignoring case and non-alphanumeric characters.
// Demonstrates usage of `strings.Builder` and `unicode` functions (`IsLetter`, `IsDigit`, `ToLower`).
func IsPalindrome(s string) bool {
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(unicode.ToLower(r))
		}
	}
	clean := builder.String()
	return clean == ReverseString(clean)
}
