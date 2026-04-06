// Package basictypes showcases Go's built-in types (slices, maps, strings, JSON). // Package declaration for basic types demonstration.
package basictypes // Defining the package name as basictypes.

import ( // Starting the import block for basic types and JSON handling.
	"encoding/json" // Importing json for encoding and decoding JSON data.
	"fmt"           // Importing fmt for formatted input/output.
	"sort"          // Importing sort for sorting slices and maps.
	"strings"       // Importing strings for string manipulation functions.
	"unicode"       // Importing unicode for character classification and transformation.
) // Closing the import block.

// ---------------------------------------------------------------------------
// 1. Slice (List) Manipulation (Standard idiomatic ways) // Section header for Slice manipulation.
// ---------------------------------------------------------------------------

// RunSliceManipulationDemo demonstrates basic slice operations. // Comment for RunSliceManipulationDemo function.
//
// For a Java developer: // Explanation targeted at Java developers.
//   - Go's `slice` is a dynamic array (like `ArrayList<T>`). // Comparing slices to Java ArrayList.
//   - Slice is a view into an underlying array. When you slice a slice, they share memory. // Explaining the view nature of slices.
//   - The zero-value of a slice is `nil` (like an uninitialized Java List). // Noting the nil zero-value.
//   - There's no formal `List` interface in the standard library. // Clarifying the lack of a List interface.
//   - Slices are passed by value, but the value is a "slice header" (pointer, length, capacity). // Explaining slice passing semantics.
func RunSliceManipulationDemo() { // RunSliceManipulationDemo function.
	fmt.Println("\n--- Slice (List) Manipulation Demo ---") // Printing the demo header.

	// a. Creation: make(type, len, cap) // Comment for slice creation.
	// Java equivalent: List<Integer> list = new ArrayList<>(10); // Comparing to Java list creation.
	nums := make([]int, 0, 5)                                                  // Creating an int slice with initial length 0 and capacity 5.
	fmt.Printf("   Initial: len=%d, cap=%d, %v\n", len(nums), cap(nums), nums) // Printing initial length and capacity.

	// b. Adding: append(slice, element...) // Comment for adding elements to a slice.
	// Java equivalent: list.add(10); // Comparing to Java list add.
	nums = append(nums, 1, 2, 3)              // Appending elements 1, 2, and 3 to the slice.
	fmt.Printf("   After append: %v\n", nums) // Printing the slice after appending.

	// c. Slicing: slice[start:end] (half-open interval [start, end)) // Comment for slicing operations.
	// Java equivalent: list.subList(1, 3); // Comparing to Java subList.
	sub := nums[1:3]                            // Creating a sub-slice from index 1 to 2 (index 3 is excluded).
	fmt.Printf("   Sub-slice [1:3]: %v\n", sub) // Printing the sub-slice.

	// d. Length and Capacity // Comment for length and capacity.
	// len() is current size, cap() is the size of the underlying array. // Explaining len vs cap.
	fmt.Printf("   Length: %d, Capacity: %d\n", len(nums), cap(nums)) // Printing current length and capacity.

	// e. Iteration: range (returns index and value) // Comment for iterating over a slice.
	// Java equivalent: for (int i=0; i < list.size(); i++) { ... } // Comparing to Java for-loop.
	fmt.Print("   Iteration: ") // Printing a label for iteration.
	for i, v := range nums {    // Iterating through the slice using range to get index and value.
		fmt.Printf("[%d]:%d ", i, v) // Printing index and value for each element.
	} // End of iteration loop.
	fmt.Println() // Printing a newline for formatting.

	// f. Copying: copy(dest, src) // Comment for copying slices.
	// Important: copy only copies the minimum of len(dest) and len(src). // Noting the limitation of the copy function.
	backup := make([]int, len(nums))             // Creating a destination slice with the same length as the source.
	copy(backup, nums)                           // Copying elements from nums to backup.
	fmt.Printf("   Copy (backup): %v\n", backup) // Printing the copied slice.
} // Closing the RunSliceManipulationDemo function.

// ---------------------------------------------------------------------------
// 2. Map Manipulation (Standard idiomatic ways) // Section header for Map manipulation.
// ---------------------------------------------------------------------------

// RunMapManipulationDemo demonstrates basic map operations. // Comment for RunMapManipulationDemo function.
//
// For a Java developer: // Explanation targeted at Java developers.
//   - Go's `map` is a hash map (like `HashMap<K, V>`). // Comparing Go maps to Java HashMaps.
//   - Iteration order is randomized to prevent dependence on it (Java doesn't guarantee order in HashMap either). // Explaining random iteration order.
//   - Maps are reference types; passing to functions modifies the original map. // Explaining map reference semantics.
//   - Accessing a non-existent key returns the zero-value (e.g., 0 for int, "" for string) instead of throwing an exception or returning null. // Explaining missing key behavior.
func RunMapManipulationDemo() { // RunMapManipulationDemo function.
	fmt.Println("\n--- Map Manipulation Demo ---") // Printing the demo header.

	// a. Creation: make(map[KeyType]ValueType) // Comment for map creation.
	// Java equivalent: Map<String, Integer> map = new HashMap<>(); // Comparing to Java map creation.
	ages := make(map[string]int) // Creating a map with string keys and int values.

	// b. Adding / Updating // Comment for adding or updating map entries.
	// Java equivalent: map.put("Alice", 30); // Comparing to Java map put.
	ages["Alice"] = 30                       // Setting the value for key "Alice" to 30.
	ages["Bob"] = 25                         // Setting the value for key "Bob" to 25.
	fmt.Printf("   Initial map: %v\n", ages) // Printing the initial map.

	// c. Existence check: the "comma ok" idiom // Comment for checking key existence.
	// Java equivalent: map.containsKey("Alice"); // Comparing to Java containsKey.
	age, ok := ages["Alice"] // Attempting to retrieve value for "Alice" and existence status.
	if ok {                  // Checking if the key was found.
		fmt.Printf("   Alice's age is %d\n", age) // Printing the age if found.
	} // End of existence check.

	// d. Deleting: delete(map, key) // Comment for deleting map entries.
	// Java equivalent: map.remove("Bob"); // Comparing to Java map remove.
	delete(ages, "Bob")                            // Deleting the entry for key "Bob" from the map.
	fmt.Printf("   After delete(Bob): %v\n", ages) // Printing the map after deletion.

	// e. Iteration (Warning: Order is random!) // Comment for map iteration.
	fmt.Print("   Iteration (order varies): ") // Printing iteration label.
	for name, age := range ages {              // Iterating through map entries using range.
		fmt.Printf("%s:%d ", name, age) // Printing name and age for each entry.
	} // End of iteration loop.
	fmt.Println() // Printing a newline.

	// f. Deterministic Iteration (Sort keys first) // Comment for deterministic map iteration.
	fmt.Print("   Deterministic iteration: ") // Printing label.
	ages["Charlie"] = 35                      // Adding key "Charlie" to the map.
	ages["Alpha"] = 20                        // Adding key "Alpha" to the map.
	keys := make([]string, 0, len(ages))      // Creating a slice to hold the map keys.
	for k := range ages {                     // Iterating through map keys to collect them.
		keys = append(keys, k) // Appending each key to the keys slice.
	} // End of collection loop.
	sort.Strings(keys)       // Sorting the collected keys alphabetically.
	for _, k := range keys { // Iterating through the sorted keys.
		fmt.Printf("%s:%d ", k, ages[k]) // Printing each key and its corresponding value.
	} // End of sorted iteration loop.
	fmt.Println() // Printing a newline.
} // Closing the RunMapManipulationDemo function.

// ---------------------------------------------------------------------------
// 3. String Operations (Standard Library) // Section header for String operations.
// ---------------------------------------------------------------------------

// RunStringOperationsDemo demonstrates common string tasks. // Comment for RunStringOperationsDemo function.
//
// For a Java developer: // Explanation targeted at Java developers.
//   - Strings in Go are immutable sequences of bytes (usually UTF-8). // Describing string nature in Go.
//   - Unlike Java's `String` which is `char[]` (UTF-16), Go strings are `byte[]`. // Contrasting with Java string internal representation.
//   - Use the `strings` package for most operations. // Recommending the strings package.
//   - String concatenation with `+` is fine for small cases, // Noting when + is acceptable.
//     but `strings.Builder` is preferred for loops (like `StringBuilder`). // Recommending strings.Builder for efficiency.
//   - Go uses backticks (“) for raw string literals (like Java 15's Text Blocks). // Comparing raw strings to Java text blocks.
func RunStringOperationsDemo() { // RunStringOperationsDemo function.
	fmt.Println("\n--- String Operations Demo ---") // Printing the demo header.

	text := "Go is a statically typed, compiled programming language." // Defining a sample string.

	// a. Basic checks // Comment for basic string status checks.
	fmt.Printf("   Contains 'typed':   %t\n", strings.Contains(text, "typed")) // Checking if string contains "typed".
	fmt.Printf("   Has prefix 'Go':    %t\n", strings.HasPrefix(text, "Go"))   // Checking if string starts with "Go".
	fmt.Printf("   Has suffix 'Java':  %t\n", strings.HasSuffix(text, "Java")) // Checking if string ends with "Java".

	// b. Manipulation // Comment for basic string transformations.
	fmt.Printf("   Upper:              %s\n", strings.ToUpper("go rocks"))              // Converting a string to uppercase.
	fmt.Printf("   Replace:            %s\n", strings.Replace(text, "Go", "Golang", 1)) // Replacing the first occurrence of "Go" with "Golang".

	// c. Splitting and Joining // Comment for splitting and joining strings.
	words := strings.Fields(text)                         // Splitting the string into a slice of words by whitespace.
	fmt.Printf("   Words count:        %d\n", len(words)) // Printing the number of words found.
	joined := strings.Join(words[:3], "-")                // Joining the first three words with hyphens.
	fmt.Printf("   Join first three:   %s\n", joined)     // Printing the joined string.

	// d. Trimming // Comment for trimming whitespace from strings.
	dirty := "   \t hello world \n  "                                    // Defining a string with excessive whitespace.
	fmt.Printf("   Trimmed:           '%s'\n", strings.TrimSpace(dirty)) // Printing the string after trimming space.

	// e. strings.Builder (Performance like StringBuilder) // Comment for strings.Builder usage.
	var builder strings.Builder // Declaring a strings.Builder instance.
	for i := 1; i <= 3; i++ {   // Loop to append multiple strings efficiently.
		builder.WriteString(fmt.Sprintf("Step %d; ", i)) // Writing a formatted string to the builder.
	} // End of loop.
	fmt.Printf("   Builder result:     %s\n", builder.String()) // Printing the final string built by the builder.

	// f. Unicode / Runes // Comment for handling Unicode characters and runes.
	// Java equivalent: String.codePointAt() // Comparing to Java's codePointAt.
	japanese := "こんにちは"                                                        // "Konnichiwa" in Japanese characters.
	fmt.Printf("   Bytes length:       %d (not characters!)\n", len(japanese)) // Printing the length of the string in bytes.
	fmt.Printf("   Runes count:        %d\n", len([]rune(japanese)))           // Printing the count of actual Unicode characters (runes).
} // Closing the RunStringOperationsDemo function.

// ---------------------------------------------------------------------------
// 4. JSON Manipulation (Standard Library) // Section header for JSON manipulation.
// ---------------------------------------------------------------------------

// Product represents a simple data model for JSON demonstration. // Comment for the Product struct.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Struct tags like `json:"id"` are similar to Jackson's `@JsonProperty("id")`. // Comparing struct tags to Jackson annotations.
// - Fields must be capitalized (Exported) for the `encoding/json` package to see them. // Noting that visibility affects serialization.
// - This is Go's way of doing meta-programming/reflection-based serialization. // Describing Go's serialization mechanism.
// - Marshalling = Serializing (Object to JSON). // Defining marshalling.
// - Unmarshalling = Deserializing (JSON to Object). // Defining unmarshalling.
type Product struct { // Defining the Product struct.
	ID    int      `json:"id"`             // Field ID mapped to JSON key "id".
	Name  string   `json:"name"`           // Field Name mapped to JSON key "name".
	Price float64  `json:"price"`          // Field Price mapped to JSON key "price".
	Tags  []string `json:"tags,omitempty"` // Field Tags mapped to "tags" and omitted if empty in JSON.
} // Closing the Product struct definition.

// RunJSONDemo demonstrates standard JSON tasks using the `encoding/json` package. // Comment for RunJSONDemo function.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Marshalling in Go is equivalent to calling `objectMapper.writeValueAsString()`. // Comparing to Jackson's Marshalling.
// - Unmarshalling is equivalent to `objectMapper.readValue()`. // Comparing to Jackson's Unmarshalling.
// - Error handling is explicit, unlike Java's `JsonProcessingException`. // Noting Go's explicit error handling style.
func RunJSONDemo() { // RunJSONDemo function.
	fmt.Println("\n--- JSON Manipulation Demo ---") // Printing the demo header.

	// a. Marshalling (Struct to JSON) // Comment for converting struct to JSON.
	p := Product{ // Initializing a Product instance.
		ID:    101,                       // Setting ID.
		Name:  "Gopher Plushie",          // Setting Name.
		Price: 19.99,                     // Setting Price.
		Tags:  []string{"toy", "mascot"}, // Setting Tags.
	} // End of initialization.
	jsonData, _ := json.MarshalIndent(p, "   ", "  ")     // Converting the Product to indented JSON bytes.
	fmt.Printf("   JSON Output:\n%s\n", string(jsonData)) // Printing the resulting JSON string.

	// b. Unmarshalling (JSON to Struct) // Comment for converting JSON back to struct.
	rawJSON := `{"id": 102, "name": "Go Mug", "price": 12.50}` // Defining a raw JSON string.
	var p2 Product                                             // Declaring a Product variable to hold the unmarshalled data.
	err := json.Unmarshal([]byte(rawJSON), &p2)                // Converting JSON bytes to the Product struct.
	if err == nil {                                            // Checking if unmarshalling succeeded.
		fmt.Printf("   Unmarshalled Struct: %+v\n", p2) // Printing the populated struct.
	} // End of success check.

	// c. Arbitrary JSON (using map[string]any) // Comment for handling dynamic JSON schemas.
	// Useful when you don't know the schema (like Java Map<String, Object>). // Explaining why dynamic mapping is useful.
	var data map[string]any                      // Declaring a map to hold dynamic JSON data.
	err = json.Unmarshal([]byte(rawJSON), &data) // Unmarshalling JSON bytes into the map.
	if err == nil {                              // Checking if unmarshalling succeeded.
		fmt.Printf("   Map representation:  %v (Name: %v)\n", data, data["name"]) // Printing the map and a specific value.
	} // End of success check.

	// d. JSON with Custom Logic (Filtering/Validation) // Comment for advanced JSON features.
	// Demonstrating 'omitempty' and ignored fields // Highlighting special JSON tag features.
	pEmpty := Product{ID: 1, Name: "Invisible"}                  // Initializing a Product with some empty fields.
	emptyJSON, _ := json.Marshal(pEmpty)                         // Marshalling the product to see omitempty in action.
	fmt.Printf("   Omitempty Tags:     %s\n", string(emptyJSON)) // Printing the JSON result.
} // Closing the RunJSONDemo function.

// RunBasicTypesDemo orchestrates the basic types demos. // Comment for RunBasicTypesDemo function.
//
// For a Java developer: // Explanation targeted at Java developers.
// - This is the "facade" for all basic types demonstrations. // Comparing to the facade pattern.
func RunBasicTypesDemo() { // RunBasicTypesDemo function.
	fmt.Println("--- Basic Types & Standard Library Demo ---") // Printing the overall section header.
	RunSliceManipulationDemo()                                 // Running the slice manipulation demo.
	RunMapManipulationDemo()                                   // Running the map manipulation demo.
	RunStringOperationsDemo()                                  // Running the string operations demo.
	RunJSONDemo()                                              // Running the JSON manipulation demo.
	fmt.Println("--- Basic Types & Standard Library End ---")  // Printing the section footer.
} // Closing the RunBasicTypesDemo function.

// ReverseString is a utility function used for demonstrating string manipulation. // Comment for the ReverseString utility.
// It uses runes to correctly handle multi-byte characters (UTF-8). // Explaining the use of runes for UTF-8 safety.
//
// For a Java developer: // Explanation targeted at Java developers.
//   - Strings in Go are UTF-8 encoded by default. // Noting Go string encoding.
//   - Iterating over a string using `range` gives you `runes` (Unicode code points), // Explaining range behavior on strings.
//     not bytes or characters. // Noting that it returns code points.
//   - `rune` is an alias for `int32`, representing a Unicode code point. // Defining what a rune is.
//   - This is necessary because Go's `string` length is in bytes, not characters. // Explaining the rationale for using runes.
func ReverseString(s string) string { // ReverseString function taking and returning a string.
	runes := []rune(s)                                    // Converting the string to a slice of runes.
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { // Two-pointer loop to reverse the rune slice.
		runes[i], runes[j] = runes[j], runes[i] // Swapping runes in-place.
	} // End of reversal loop.
	return string(runes) // Converting the reversed rune slice back to a string and returning it.
} // Closing the ReverseString function.

// IsPalindrome checks if a string is a palindrome, ignoring case and non-letters. // Comment for IsPalindrome function.
//
// For a Java developer: // Explanation targeted at Java developers.
// - Demonstrates usage of `strings.Builder` and the `unicode` package. // Highlighting key standard library usage.
// - `strings.Builder` is the Go equivalent of `StringBuilder`. // Comparing to Java StringBuilder.
// - Shows how to iterate over runes in a string. // Noting the demonstration of rune iteration.
func IsPalindrome(s string) bool { // IsPalindrome function returning a boolean.
	var builder strings.Builder // Declaring a strings.Builder for normalized string construction.
	for _, r := range s {       // Iterating through each rune of the input string.
		if unicode.IsLetter(r) || unicode.IsDigit(r) { // Filtering for letters and digits.
			builder.WriteRune(unicode.ToLower(r)) // Writing the lowercase version of the rune to the builder.
		} // End of filtering logic.
	} // End of character collection loop.
	clean := builder.String()            // Getting the normalized string from the builder.
	return clean == ReverseString(clean) // Checking if the cleaned string is equal to its reverse.
} // Closing the IsPalindrome function.
