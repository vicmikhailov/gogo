package interfaces

import (
	"math"
	"testing"
)

/**
 * ===========================================================================
 * Interface and Struct Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Go tests are always functions starting with `Test` and taking `(t *testing.T)`.
 * - No `@Test` annotation is required.
 * - Assertions are manual: `if condition { t.Errorf(...) }`.
 * - `t.Errorf` does NOT stop execution (like a soft assertion). Use `t.Fatalf`
 *   if you want to stop immediately.
 */

func TestRectangle(t *testing.T) {
	r := Rectangle{Width: 10, Height: 5}
	if r.Area() != 50 {
		t.Errorf("Expected Area 50, got %.2f", r.Area())
	}
	if r.Perimeter() != 30 {
		t.Errorf("Expected Perimeter 30, got %.2f", r.Perimeter())
	}
}

func TestCircle(t *testing.T) {
	c := Circle{Radius: 7}
	expectedArea := math.Pi * 49
	if math.Abs(c.Area()-expectedArea) > 1e-9 {
		t.Errorf("Expected Area %.2f, got %.2f", expectedArea, c.Area())
	}
	expectedPerimeter := 2 * math.Pi * 7
	if math.Abs(c.Perimeter()-expectedPerimeter) > 1e-9 {
		t.Errorf("Expected Perimeter %.2f, got %.2f", expectedPerimeter, c.Perimeter())
	}
}

// TestShapeInterface showcases table-driven tests for an interface.
func TestShapeInterface(t *testing.T) {
	tests := []struct {
		name      string
		shape     Shape
		wantArea  float64
		wantPerim float64
	}{
		{"Rectangle", Rectangle{10, 5}, 50, 30},
		{"Circle", Circle{7}, math.Pi * 49, 2 * math.Pi * 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if math.Abs(tt.shape.Area()-tt.wantArea) > 1e-9 {
				t.Errorf("Area() = %.2f; want %.2f", tt.shape.Area(), tt.wantArea)
			}
			if math.Abs(tt.shape.Perimeter()-tt.wantPerim) > 1e-9 {
				t.Errorf("Perimeter() = %.2f; want %.2f", tt.shape.Perimeter(), tt.wantPerim)
			}
		})
	}
}
