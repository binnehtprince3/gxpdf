// Package table provides domain entities for PDF table extraction.
package table

import "fmt"

// Rectangle represents a rectangular bounding box.
//
// This is a value object in DDD - it's immutable and compared by value.
// Coordinates follow PDF convention: origin at bottom-left, Y increases upward.
type Rectangle struct {
	X      float64 // Bottom-left X coordinate
	Y      float64 // Bottom-left Y coordinate
	Width  float64 // Width
	Height float64 // Height
}

// NewRectangle creates a new Rectangle.
func NewRectangle(x, y, width, height float64) Rectangle {
	return Rectangle{X: x, Y: y, Width: width, Height: height}
}

// Right returns the X coordinate of the right edge.
func (r Rectangle) Right() float64 {
	return r.X + r.Width
}

// Top returns the Y coordinate of the top edge.
func (r Rectangle) Top() float64 {
	return r.Y + r.Height
}

// Bottom returns the Y coordinate of the bottom edge.
func (r Rectangle) Bottom() float64 {
	return r.Y
}

// Left returns the X coordinate of the left edge.
func (r Rectangle) Left() float64 {
	return r.X
}

// Contains checks if a point (x, y) is inside the rectangle.
func (r Rectangle) Contains(x, y float64) bool {
	return x >= r.X && x <= r.Right() && y >= r.Y && y <= r.Top()
}

// String returns a string representation of the rectangle.
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle{x=%.2f, y=%.2f, w=%.2f, h=%.2f}", r.X, r.Y, r.Width, r.Height)
}
