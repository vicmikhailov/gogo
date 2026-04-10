// Package syntax showcases Go's core syntax and how it differs from Java.
package syntax

import (
	"fmt"
)

// ---------------------------------------------------------------------------
// 1. Package-Level Visibility (Public vs. Private)
// ---------------------------------------------------------------------------

// PublicVariable is accessible from other packages because it starts with an uppercase letter.
// Java equivalent: `public static int PUBLIC_VARIABLE = 10;`
var PublicVariable = 10

// privateVariable is ONLY accessible within the `syntax` package.
// Java equivalent: `private static int PRIVATE_VARIABLE = 20;`
var privateVariable = 20

// ---------------------------------------------------------------------------
// 2. Constants and iota (The "Go Enum")
// ---------------------------------------------------------------------------

// Status represents a custom type for iota-based enums.
type Status int

// Java equivalent: `enum Status { PENDING, RUNNING, COMPLETED }`
const (
	Pending   Status = iota // 0
	Running                 // 1
	Completed               // 2
)

// ---------------------------------------------------------------------------
// 3. Named Return Values
// ---------------------------------------------------------------------------

// CalculateDimensions returns both the area and perimeter of a square.
// For a Java developer:
// - Go functions can return multiple values.
// - Named return values are initialized to their zero-values.
// - You can use `return` without arguments (a "naked" return) to return the current values.
func CalculateDimensions(side float64) (area float64, perimeter float64) {
	area = side * side
	perimeter = 4 * side
	return // Returns the current values of area and perimeter.
}

// ---------------------------------------------------------------------------
// 4. Pointers vs Values (The "Big One")
// ---------------------------------------------------------------------------

// Point is a simple struct.
type Point struct {
	X, Y int
}

// UpdateValue takes a copy of the struct (Pass-by-value).
// The original struct is NOT modified.
func UpdateValue(p Point) {
	p.X = 100
}

// UpdatePointer takes a pointer to the struct (Pass-by-reference-like behavior).
// The original struct IS modified.
func UpdatePointer(p *Point) {
	p.X = 100
}

// ---------------------------------------------------------------------------
// 5. Deferred Execution (defer)
// ---------------------------------------------------------------------------

// DeferExample showcases the `defer` keyword.
// For a Java developer:
// - `defer` is like a `finally` block but more flexible.
// - It schedules a function call to run just before the surrounding function returns.
// - Defers are executed in Last-In-First-Out (LIFO) order.
func DeferExample() {
	defer fmt.Println("      (This runs LAST, like a finally block)")
	fmt.Println("      (This runs first)")
}

// ---------------------------------------------------------------------------
// 6. RunSyntaxDemo
// ---------------------------------------------------------------------------

// RunSyntaxDemo showcases the fundamental syntax differences for Java developers.
func RunSyntaxDemo() {
	fmt.Println("--- Go Syntax for Java Developers ---")

	// a. Variable Declarations
	// Java: `int x = 10;`
	// Go: Multiple ways!
	var x int = 10 // Explicit type
	var y = 20     // Type inference
	z := 30        // Short-hand declaration (only inside functions)
	fmt.Printf("   1. Variables: x=%d, y=%d, z=%d\n", x, y, z)

	// b. Multiple Return Values
	area, peri := CalculateDimensions(5.0)
	fmt.Printf("   2. Multiple Returns: Area=%.2f, Perimeter=%.2f\n", area, peri)

	// c. Blank Identifier (_)
	// If you don't need a value, you MUST use `_` to discard it.
	// In Java, you just don't assign it. In Go, an unused variable is a COMPILE error!
	a, _ := CalculateDimensions(10.0)
	fmt.Printf("   3. Blank Identifier: Discarded perimeter, Area=%.2f\n", a)

	// d. For is the only loop!
	// Java has `for`, `while`, `do-while`. Go only has `for`.
	fmt.Print("   4. Loops (For is everything): ")
	// While-like behavior:
	count := 0
	for count < 3 {
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println("(The only loop keyword in Go)")

	// e. Switch (No fallthrough by default)
	// Java requires `break` for every case. Go doesn't!
	fmt.Print("   5. Switch (Safe by default): ")
	status := Running
	switch status {
	case Pending:
		fmt.Print("Pending")
	case Running:
		fmt.Print("Running (No break needed!)")
	case Completed:
		fmt.Print("Completed")
	}
	fmt.Println()

	// f. Pointers vs Values
	p := Point{X: 1, Y: 1}
	UpdateValue(p)
	fmt.Printf("   6. Pointers: After UpdateValue: X=%d (No change)\n", p.X)
	UpdatePointer(&p)
	fmt.Printf("      After UpdatePointer: X=%d (Changed!)\n", p.X)

	// g. Deferred Execution
	fmt.Println("   7. Defer (The Go 'finally'):")
	DeferExample()

	// h. Type Conversion (No implicit casts!)
	// Java: `double d = 10;` (Implicit cast from int)
	// Go: Requires explicit conversion.
	var integer int = 42
	var f float64 = float64(integer)
	fmt.Printf("   8. Type Conversion: int %d -> float %.2f (Explicit!)\n", integer, f)

	// i. Zero-Values (Zero Initialization)
	// In Java, local variables must be initialized. In Go, they are ALWAYS zeroed.
	var uninitializedInt int
	var uninitializedBool bool
	var uninitializedString string
	fmt.Printf("   9. Zero-Values: int=%d, bool=%v, string='%s'\n", uninitializedInt, uninitializedBool, uninitializedString)

	fmt.Println("--- Go Syntax Demo End ---")
}
