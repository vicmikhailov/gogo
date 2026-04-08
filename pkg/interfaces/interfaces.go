// Package interfaces showcases Go's implicit interface system and polymorphism.
package interfaces

import (
	"fmt"
	"math"
)

// Shape is an interface that defines the common behavior for shapes.
//
// For a Java developer:
//   - Go interfaces are "implicit". There is no `implements` keyword.
//   - If a `struct` has all the methods defined in an `interface`, it
//     automatically satisfies that interface. This is often called "duck typing"
//     or "structural typing".
//   - This allows you to define interfaces for types you don't even own (retroactive implementation)!
//   - Java-ism to avoid: Defining an interface in the same package as the implementation.
//   - Go idiom: "Accept interfaces, return structs". Define interfaces in the CONSUMER package,
//     specifying only the methods that the consumer actually needs.
type Shape interface {
	// Area returns the surface area of the shape.
	Area() float64
	// Perimeter returns the distance around the shape.
	Perimeter() float64
}

// Rectangle is a struct that represents a rectangle.
//
// For a Java developer:
//   - A `struct` is similar to a POJO or a Record. It only holds data.
//   - It has no constructors; you use "composite literals" like `Rectangle{Width: 10}`.
//   - Methods are defined OUTSIDE the struct, but "bound" to it via a "receiver".
type Rectangle struct {
	Width, Height float64
}

// Area returns the area of the rectangle.
// Java equivalent: `public double area()`
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter returns the perimeter of the rectangle.
// Java equivalent: `public double perimeter()`
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// It satisfies Shape implicitly because it implements Area and Perimeter.
//
// For a Java developer:
// - There is no `implements Shape` here, but Go sees the methods match.
type Circle struct {
	Radius float64
}

// Area returns the area of the circle.
// Java equivalent: `public double area()`
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter returns the circumference of the circle.
// Java equivalent: `public double perimeter()`
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// PrintShapeDetails accepts any Shape and prints its details.
// Java equivalent: `public void printShapeDetails(Shape s)`
func PrintShapeDetails(s Shape) {
	// %T prints the type of the variable.
	fmt.Printf("   Type: %T, Area: %.2f, Perimeter: %.2f\n", s, s.Area(), s.Perimeter())
}

// RunInterfacesDemo showcases Go's interface system and polymorphism.
func RunInterfacesDemo() {
	fmt.Println("--- Interfaces & Structs Demo ---")

	r := Rectangle{Width: 10, Height: 5}
	c := Circle{Radius: 7}

	// 1. Polymorphism with interfaces
	// Java equivalent: List<Shape> shapes = Arrays.asList(r, c);
	fmt.Println("1. Polymorphism with interfaces:")
	shapes := []Shape{r, c}
	for _, s := range shapes {
		PrintShapeDetails(s)
	}

	// 2. Type assertion
	// Java equivalent: `if (s instanceof Circle) { Circle circle = (Circle) s; ... }`
	fmt.Println("2. Type assertion:")
	var s Shape = Circle{Radius: 10}
	if circle, ok := s.(Circle); ok { // 'ok' is true if the assertion succeeded
		fmt.Printf("   Successfully asserted as Circle with Radius: %.2f\n", circle.Radius)
	}

	// 3. Type switch
	// Java equivalent: Modern Java `switch` with pattern matching.
	fmt.Println("3. Type switch:")
	for _, shape := range shapes {
		switch v := shape.(type) {
		case Rectangle:
			fmt.Printf("   Found a Rectangle with width %.2f and height %.2f\n", v.Width, v.Height)
		case Circle:
			fmt.Printf("   Found a Circle with radius %.2f\n", v.Radius)
		}
	}

	fmt.Println("--- Interfaces & Structs Demo End ---")
}
