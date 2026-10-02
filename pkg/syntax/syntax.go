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

// privateVariable is ONLY accessible within the `syntax` package because it starts with a lowercase letter.
// Java equivalent: `private static int PRIVATE_VARIABLE = 20;`
var privateVariable = 20

// ---------------------------------------------------------------------------
// 2. Constants and iota (The "Go Enum")
// ---------------------------------------------------------------------------

// Status represents a custom type for iota-based enums.
type Status int

// Java equivalent: `enum Status { PENDING, RUNNING, COMPLETED }`
// `iota` auto-increments constants within a const block starting at 0:
// Pending = 0, Running = 1, Completed = 2.
const (
	Pending Status = iota
	Running
	Completed
)

// ---------------------------------------------------------------------------
// 3. Named Return Values
// ---------------------------------------------------------------------------

// CalculateDimensions returns both the area and perimeter of a square.
//
// For a Java developer:
//   - Go functions can return multiple values.
//   - Named return values (`area`, `perimeter`) are pre-declared and initialized to zero values.
//   - A "naked" `return` returns the current values of the named return variables.
func CalculateDimensions(side float64) (area float64, perimeter float64) {
	area = side * side
	perimeter = 4 * side
	return
}

// ---------------------------------------------------------------------------
// 4. Pointers vs Values (The "Big One")
// ---------------------------------------------------------------------------

// Point is a simple 2D coordinate struct.
type Point struct {
	X, Y int
}

// UpdateValue receives a copy of the Point struct (pass-by-value).
// The caller's original struct remains unchanged.
func UpdateValue(p Point) {
	p.X = 100
}

// UpdatePointer receives a pointer to the Point struct (*Point).
// The caller's original struct is modified in-place.
func UpdatePointer(p *Point) {
	p.X = 100
}

// ---------------------------------------------------------------------------
// 5. Deferred Execution (defer)
// ---------------------------------------------------------------------------

// DeferExample showcases the `defer` keyword.
//
// For a Java developer:
//   - `defer` schedules a function call to execute right before the surrounding function returns.
//   - Similar to a `finally` block, but scoped to the surrounding function rather than a block.
//   - Multiple defers are executed in Last-In-First-Out (LIFO) order.
func DeferExample() {
	defer fmt.Println("      (This runs LAST, like a finally block)")
	fmt.Println("      (This runs first)")
}

// ---------------------------------------------------------------------------
// 6. RunSyntaxDemo
// ---------------------------------------------------------------------------

// RunSyntaxDemo showcases fundamental syntax differences for Java developers.
func RunSyntaxDemo() {
	fmt.Println("--- Go Syntax for Java Developers ---")

	// a. Variable Declarations
	// Java: `int x = 10;`
	// Go provides multiple declaration styles:
	//   - Explicit type: var x int = 10
	//   - Type inference: var y = 20
	//   - Short-hand (inside functions only): z := 30
	var x int = 10
	var y = 20
	z := 30
	fmt.Printf("   1. Variables: x=%d, y=%d, z=%d\n", x, y, z)

	// b. Multiple Return Values
	// Java requires a Tuple, Record, or custom wrapper class to return multiple values.
	area, peri := CalculateDimensions(5.0)
	fmt.Printf("   2. Multiple Returns: Area=%.2f, Perimeter=%.2f\n", area, peri)

	// c. Blank Identifier (_)
	// Unused variables are compile-time errors in Go.
	// Use `_` to intentionally discard values you do not need.
	a, _ := CalculateDimensions(10.0)
	fmt.Printf("   3. Blank Identifier: Discarded perimeter, Area=%.2f\n", a)

	// d. Loops (For is the only loop!)
	// Java has `for`, `while`, and `do-while`. Go only has the `for` keyword.
	// A condition-only for loop behaves identically to while.
	fmt.Print("   4. Loops (For is everything): ")
	count := 0
	for count < 3 {
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println("(The only loop keyword in Go)")

	// e. Switch (No fallthrough by default)
	// Unlike Java switch statements that require explicit `break`, Go cases break automatically.
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
	// Java permits implicit widening casts (e.g. double d = 10;).
	// Go requires explicit conversion: float64(integer).
	var integer int = 42
	var f float64 = float64(integer)
	fmt.Printf("   8. Type Conversion: int %d -> float %.2f (Explicit!)\n", integer, f)

	// i. Zero-Values (Zero Initialization)
	// In Java, uninitialized local variables cause a compiler error.
	// In Go, variables are ALWAYS guaranteed to be initialized to their zero-value (0, false, "").
	var uninitializedInt int
	var uninitializedBool bool
	var uninitializedString string
	fmt.Printf("   9. Zero-Values: int=%d, bool=%v, string='%s'\n", uninitializedInt, uninitializedBool, uninitializedString)

	fmt.Println("--- Go Syntax Demo End ---")
}
