package syntax

import (
	"testing"
)

func TestCalculateDimensions(t *testing.T) {
	area, peri := CalculateDimensions(5.0)
	if area != 25.0 {
		t.Errorf("Expected area 25.0, got %f", area)
	}
	if peri != 20.0 {
		t.Errorf("Expected perimeter 20.0, got %f", peri)
	}
}

func TestUpdatePointer(t *testing.T) {
	p := Point{X: 1, Y: 1}
	UpdatePointer(&p)
	if p.X != 100 {
		t.Errorf("Expected X to be 100, got %d", p.X)
	}
}

func TestUpdateValue(t *testing.T) {
	p := Point{X: 1, Y: 1}
	UpdateValue(p)
	if p.X != 1 {
		t.Errorf("Expected X to remain 1, got %d", p.X)
	}
}
